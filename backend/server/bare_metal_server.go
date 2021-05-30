package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//BareMetalServer 裸机环境
type BareMetalServer struct {
	NacosNameClient naming_client.INamingClient
	NacosIP         string
	NacosPort       uint64
	ServiceName     string
	MonitorServices []string
	ServiceGroup    string
}

//NewMonitorServer 创建新的BareMetalServer
func NewBareMetalServer(config vo.NacosClientParam, nacosIP string, nacosPort uint64,
	monitorServices []string, group string) (*BareMetalServer, error) {
	var server BareMetalServer
	server.ServiceName = "gomonitor-server"
	server.MonitorServices = monitorServices
	server.ServiceGroup = group
	server.NacosIP = nacosIP
	server.NacosPort = nacosPort
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
		Ip:          server.NacosIP,
		Port:        server.NacosPort,
		ServiceName: server.ServiceName,
		Enable:      false,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"info": "hello"},
	})
	//开启监控服务
	for _, service := range server.MonitorServices {
		log.Info("start monitor service: " + service)
		server.NacosNameClient.Subscribe(&vo.SubscribeParam{
			ServiceName:       service,
			GroupName:         server.ServiceGroup,
			SubscribeCallback: callback,
		})
	}
}

//callback
//TODO: check valid
func callback(services []model.SubscribeService, err error) {
	for _, service := range services {
		if err != nil {
			log.Error(err)
			continue
		}
		pid := service.Metadata["pid"]
		log.Info("start monitor prog on ip " + service.Ip + " for " + pid)
		startMonitorProg(service.Ip, pid)
	}
}

//在对应ip启动monitor进程
func startMonitorProg(ip string, pid string) {

}
