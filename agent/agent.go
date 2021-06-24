package main

import (
	"context"
	"kube-gomonitor/agent/internal"
	"kube-gomonitor/agent/packet"
	"kube-gomonitor/agent/process"
	"kube-gomonitor/agent/report"
	"kube-gomonitor/pkg"
	"os"
	"os/signal"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	gpprocess "github.com/shirou/gopsutil/process"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"

	"syscall"

	flag "github.com/spf13/pflag"
)

var (
	flags         internal.CmdFlags
	reporter      report.Reporter
	namingClients []naming_client.INamingClient
)

func init() {
	flag.BoolVarP(&flags.Help, "help", "h", false, "gomonitor help")

	flag.Uint64Var(&flags.NacosPort, "nacosPort", 8848, "nacos server port")
	flag.StringSliceVar(&flags.NacosIPs, "nacosIPs", nil, "nacos server ips")

	flag.StringVar(&flags.MonitorIP, "monitorIP", "", "ip to be monitored")
	flag.StringVar(&flags.MonitorService, "monitorService", "", "service name to be monitored")
	flag.Int32Var(&flags.MonitorPid, "monitorPid", -1, "pid to be monitored")
	flag.Int64Var(&flags.MonitorInterval, "monitorInterval", 1, "interval seconds to send monitor info")

	flag.StringVar(&flags.DBUrl, "dburl", "", "data base url")
	flag.StringVar(&flags.Bucket, "bucket", "", "data base bucket for influxdb")
	flag.StringVar(&flags.Organization, "organization", "", "data base org for influxdb")
	flag.StringVar(&flags.Token, "token", "", "data base token for influxdb")
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	flag.Parse()

	var monitorpid int32

	_, err := pkg.GetKubeConfig()
	if err == rest.ErrNotInCluster {
		//running on bare metal
		log.Info("running on bare metal")
		if flags.Help || flags.MonitorIP == "" || flags.MonitorPid == -1 ||
			flags.MonitorService == "" || flags.Bucket == "" || flags.DBUrl == "" ||
			flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			os.Exit(0)
		}

		monitorpid = flags.MonitorPid
		// 在nacos注册自己
		if flags.NacosIPs != nil {
			namingClients, err := pkg.GetNacosNamingClients(flags.NacosIPs, flags.NacosPort, []string{"DEFAULT_GROUP"})
			if err != nil {
				log.Error(err)
				return
			}
			namingClients[0].RegisterInstance(vo.RegisterInstanceParam{
				Ip:          flags.MonitorIP,
				Port:        0,
				ServiceName: "kube-gomonitor agent",
				Enable:      false,
				Healthy:     true,
				Ephemeral:   true,
				Metadata:    map[string]string{"kube-gomonitor-agent": "hello"},
			})
		}

	} else {
		//running in k8s
		log.Info("running on kubernetes")
		if flags.Help || flags.Bucket == "" || flags.DBUrl == "" ||
			flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			os.Exit(0)
		}

		//get pid need to be monitored
		pids, err := gpprocess.Pids()
		if err != nil {
			log.Error(err)
			return
		}
		for _, v := range pids {
			if v != int32(os.Getpid()) {
				monitorpid = v
				break
			}
		}
	}

	go packet.NetSniff(ctx, flags.MonitorIP)

	processInfo, err := process.NewProcessInfo(monitorpid, flags.MonitorService, flags.MonitorIP)
	if err != nil {
		log.Error(err)
		return
	}
	reporter = report.NewInfluxDBReporter(&flags, processInfo)
	defer reporter.Close()
	reporter.Start(ctx)

	<-stopCh
	cancel()
	if err == rest.ErrNotInCluster {
		namingClients[0].DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          flags.MonitorIP,
			Port:        0,
			ServiceName: "kube-gomonitor backend",
			Ephemeral:   true,
		})
	}
	time.Sleep(5 * time.Second)
}
