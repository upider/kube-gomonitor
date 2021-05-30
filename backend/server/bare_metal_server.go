package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var agentImage string = "1445277435/gomonitor-agent:v0.0.1"

type ServerFlags struct {
	NacosIP             string
	NacosPort           uint64
	NamespaceId         string
	ServiceName         string
	ServerIP            string
	MonitorServices     []string
	MonitorServiceGroup string
	DBUrl               string
	Organization        string
	Bucket              string
	Token               string
}

//BareMetalServer 裸机环境
type BareMetalServer struct {
	NacosNameClient naming_client.INamingClient
	flags           *ServerFlags
}

//NewMonitorServer 创建新的BareMetalServer
func NewBareMetalServer(config vo.NacosClientParam, flags *ServerFlags) (*BareMetalServer, error) {
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
		log.Info("start monitor prog on ip " + ip + " for " + pid)
		go func() {
			server.startMonitorProg(ip, pid, serviceName)
		}()
	}
}

//在对应ip启动monitor进程
func (server *BareMetalServer) startMonitorProg(ip string, pid string, serviceName string) {
	cli, err := dockerClient.NewClient(ip, "1.13.1", nil, nil)

	if err != nil {
		log.Error(err)
		return
	}

	ctx := context.Background()
	//容器网络设置为hostnetwork
	//并设置环境变量
	envs := []string{"MONITOR_PID=" + ip, "MONITOR_SERVICE=" + serviceName, "MONITOR_PID=" + pid,
		"REPORT_DBURL=" + server.flags.DBUrl, "REPORT_DBBUCKET=" + server.flags.Bucket,
		"REPORT_DBORG=" + server.flags.Organization, "REPORT_DBTOKEN=" + server.flags.Token}

	cli.ContainerCreate(ctx, &container.Config{
		Image:      agentImage,
		User:       "root",
		WorkingDir: "/root",
		Env:        envs,
	}, &container.HostConfig{
		NetworkMode: "host",
	}, nil, nil, "gomonitor")

	err = cli.Close()
	if err != nil {
		log.Error(err)
		return
	}
}
