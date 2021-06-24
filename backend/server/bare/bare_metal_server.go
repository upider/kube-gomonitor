package bare

import (
	"context"
	"fmt"
	"kube-gomonitor/backend/server"
	"kube-gomonitor/pkg"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var agentContainerName = "kubeGomonitorAgent"

//BareMetalServer 裸机环境
type BareMetalServer struct {
	NacosNameClients []naming_client.INamingClient
	flags            *server.ServerFlags
}

//NewMonitorServer 创建新的BareMetalServer
func NewBareMetalServer(flags *server.ServerFlags) (*BareMetalServer, error) {
	server := &BareMetalServer{}

	nacosNamingClients, err := pkg.GetNacosNamingClients(flags.NacosIPs, flags.NacosPort, flags.Namespaces)

	if err != nil {
		log.Errorf("get nacosNamingClients error: %v", err)
	}
	server.NacosNameClients = nacosNamingClients
	server.flags = flags
	return server, nil
}

//Start 启动服务
func (server *BareMetalServer) Start() {
	ip, err := pkg.GetLocalIP()
	if err != nil {
		log.Errorf("get ip error: %v", err)
	}
	//在nacos注册自己
	server.NacosNameClients[0].RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        0,
		ServiceName: "kube-gomonitor backend",
		Enable:      false,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"kube-gomonitor-backend": "hello"},
	})

	//开启监控服务
	for _, nacosNameClient := range server.NacosNameClients {
		for _, group := range server.flags.MonitorServiceGroups {
			for _, service := range server.flags.MonitorServices {
				log.Info("start monitoring services: " + service)
				nacosNameClient.Subscribe(&vo.SubscribeParam{
					ServiceName:       service,
					GroupName:         group,
					SubscribeCallback: server.callback,
				})
			}
		}
	}

	// TODO: 监控kube-gomonitor-agent
}

//callback
//TODO: check valid, valid暂时不可用
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
			log.Infof("start monitor agent on %s for service %s", ip, serviceName)
			server.startMonitorProg(ip, pid, serviceName)
		}()

		// if service.Valid {
		// 	go func() {
		// 		log.Infof("start monitor agent on %s for service  ", ip, serviceName)
		// 		server.startMonitorProg(ip, pid, serviceName)
		// 	}()
		// } else {
		// 	log.Infof("service %s on %s has stopped.", service.ServiceName, service.Ip)
		// }
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
		log.Errorf("new client error: %v", err.Error())
		return
	}

	ctx := context.Background()
	//容器网络设置为hostnetwork
	//并设置环境变量
	envs := []string{"MONITOR_IP=" + ip, "MONITOR_SERVICE=" + serviceName, "MONITOR_PID=" + pid,
		"REPORT_DBURL=" + server.flags.DBUrl, "REPORT_DBBUCKET=" + server.flags.Bucket,
		"REPORT_DBORG=" + server.flags.Organization, "REPORT_DBTOKEN=" + server.flags.Token,
		"NACOS_IP=" + server.flags.NacosIPs[0], "NACOS_PORT=" + strconv.FormatUint(server.flags.NacosPort, 10),
		"MONITOR_INTERVAL=" + strconv.FormatUint(server.flags.Interval, 10)}

	out, err := cli.ImagePull(ctx, server.flags.AgentImage, types.ImagePullOptions{})
	if err != nil {
		log.Errorf("container pull error: %v", err.Error())
		return
	}
	defer out.Close()

	containerName := fmt.Sprintf("%s-%s", agentContainerName, pid)
	err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		log.Warnf("container remove warn: %s", err.Error())
	}

	container, err := cli.ContainerCreate(ctx, &container.Config{
		Image: server.flags.AgentImage,
		User:  "root",
		Env:   envs,
		Tty:   true,
	}, &container.HostConfig{
		NetworkMode: "host", //use host PidMode and host NetworkMode, we can get info on the host node
		Privileged:  true,
		PidMode:     "host",
		// AutoRemove:  true, TODO: 容器退出后自动删除, 代码完成后使用该特性
	}, nil, nil, containerName)

	if err != nil {
		log.Errorf("container create error: %s", err.Error())
		return
	}

	err = cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Error("container start error: " + err.Error())
		return
	}

	err = cli.Close()
	if err != nil {
		log.Error(err)
		return
	}
}
