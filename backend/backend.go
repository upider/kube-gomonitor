package main

import (
	"kube-gomonitor/backend/server"
	"kube-gomonitor/backend/server/bare"
	"kube-gomonitor/backend/server/k8s"
	"kube-gomonitor/pkg"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/rest"
)

var (
	flags server.ServerFlags
	//server
	monitorServer server.MonitorServer
)

func init() {
	flag.Uint64Var(&flags.NacosPort, "nacosPort", 8848, "nacos server port")
	flag.StringSliceVar(&flags.NacosIPs, "nacosIPs", nil, "nacos server ips")

	flag.StringVar(&flags.AgentImage, "agentImage", "1445277435/kube-gomonitor-agent:v0.0.1", "monitor agent docker image")
	flag.StringSliceVar(&flags.MonitorServiceGroups, "groups", []string{"DEFAULT_GROUP"}, "monitor service groups")
	flag.StringSliceVar(&flags.MonitorServices, "services", nil, "monitor service names")
	flag.StringSliceVar(&flags.Namespaces, "namespaces", []string{"public"}, "nacos namespace ids")
	flag.Uint64Var(&flags.Interval, "interval", 3, "agent monitor interval")

	flag.StringVar(&flags.DBUrl, "dburl", "", "data base url for influxdb")
	flag.StringVar(&flags.Bucket, "bucket", "", "data base bucket for influxdb")
	flag.StringVar(&flags.Organization, "organization", "", "data base org for influxdb")
	flag.StringVar(&flags.Token, "token", "", "data base token for influxdb")
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	_, err := pkg.GetKubeConfig()

	if err == rest.ErrNotInCluster {
		flag.Parse()
		if flags.NacosIPs == nil || flags.MonitorServices == nil || flags.DBUrl == "" ||
			flags.Bucket == "" || flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			return
		}

		log.Info("running on bare metal")

		monitorServer, err = bare.NewBareMetalServer(&flags)

		if err != nil {
			log.Errorln(err)
			return
		}

		monitorServer.Start()

		<-signalChan
		time.Sleep(5 * time.Second)

	} else {
		flag.Parse()
		if flags.DBUrl == "" || flags.Bucket == "" ||
			flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			return
		}

		log.Info("running on kubernetes")
		monitorServer = k8s.NewKServer(&flags)
		monitorServer.Start()

		<-signalChan
		time.Sleep(5 * time.Second)
	}
}
