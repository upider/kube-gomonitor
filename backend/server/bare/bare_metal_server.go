package bare

import (
	"context"
	"gomonitor/backend/server"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var agentImage string = "1445277435/gomonitor-agent:v0.0.1"

//BareMetalServer 裸机环境
type BareMetalServer struct {
	NacosNameClient naming_client.INamingClient
	flags           *server.ServerFlags
}

//NewMonitorServer 创建新的BareMetalServer
func NewBareMetalServer(config vo.NacosClientParam, flags *server.ServerFlags) (*BareMetalServer, error) {
	var server BareMetalServer
	server.flags = flags
	nacosServer, err := clients.NewNamingClient(config)
	if err != nil {
		return nil, err
	}
	server.NacosNameClient = nacosServer
	return &server, nil
}

//Start 启动服务
func (server *BareMetalServer) Start(ctx context.Context) {
	//在nacos注册自己
	server.NacosNameClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          server.flags.NacosIP,
		Port:        server.flags.NacosPort,
		ServiceName: "gomonitor-server",
		Enable:      false,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"info": "hello"},
	})
	//开启监控服务
	for _, service := range server.flags.MonitorServices {
		log.Info("start monitoring service: " + service)
		server.NacosNameClient.Subscribe(&vo.SubscribeParam{
			ServiceName:       service,
			GroupName:         server.flags.MonitorServiceGroup,
			SubscribeCallback: server.callback,
		})
	}
}

//callback
//TODO: check valid
func (server *BareMetalServer) callback(services []model.SubscribeService, err error) {
	for _, service := range services {
		if err != nil {
			log.Error(err)
			continue
		}
		pid := service.Metadata["pid"]
		ip := service.Ip
		serviceName := service.ServiceName
		go func() {
			log.Info("start monitor prog on ip " + ip + " for " + pid)
			server.startMonitorProg(ip, pid, serviceName)
		}()
	}
}

//在对应ip启动monitor进程
func (server *BareMetalServer) startMonitorProg(ip string, pid string, serviceName string) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}
	httpClient := &http.Client{Transport: tr}
	cli, err := dockerClient.NewClient("http://"+ip+":2375", "1.41", httpClient, nil)

	if err != nil {
		log.Error("new client error: " + err.Error())
		return
	}

	ctx := context.Background()
	//容器网络设置为hostnetwork
	//并设置环境变量
	envs := []string{"MONITOR_IP=" + ip, "MONITOR_SERVICE=" + serviceName, "MONITOR_PID=" + pid,
		"REPORT_DBURL=" + server.flags.DBUrl, "REPORT_DBBUCKET=" + server.flags.Bucket,
		"REPORT_DBORG=" + server.flags.Organization, "REPORT_DBTOKEN=" + server.flags.Token}

	out, err := cli.ImagePull(ctx, agentImage, types.ImagePullOptions{})
	if err != nil {
		log.Error("container pull error: " + err.Error())
		return
	}
	defer out.Close()

	err = cli.ContainerRemove(ctx, "gomonitor", types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Warn("container remove error: " + err.Error())
		// return
	}

	container, err := cli.ContainerCreate(ctx, &container.Config{
		Image: agentImage,
		User:  "root",
		Env:   envs,
		Tty:   true,
	}, &container.HostConfig{
		NetworkMode: "host", //use host PidMode and host NetworkMode, we can get info on the host node
		Privileged:  true,
		PidMode:     "host",
		// AutoRemove:  true, TODO: 容器退出后自动删除, 代码完成后使用该特性
	}, nil, nil, "gomonitor")

	if err != nil {
		log.Error("container create error: " + err.Error())
		return
	}

	err = cli.Close()
	if err != nil {
		log.Error(err)
		return
	}

	err = cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Error("container start error: " + err.Error())
		return
	}
}
