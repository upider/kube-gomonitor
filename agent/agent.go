package main

import (
	"context"
	"gomonitor/agent/packet"
	"gomonitor/agent/process"
	"gomonitor/agent/report"
	"gomonitor/utils"
	"os"
	"os/signal"
	"time"

	gpprocess "github.com/shirou/gopsutil/process"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"

	"syscall"

	flag "github.com/spf13/pflag"
)

type cmdFlags struct {
	Help            bool
	MonitorIP       string
	MonitorPid      int32
	MonitorService  string
	MonitorInterval int64
	DBUrl           string
	Organization    string
	Bucket          string
	Token           string
}

var (
	flags    cmdFlags
	reporter report.Reporter
)

func init() {
	flag.BoolVarP(&flags.Help, "help", "h", false, "gomonitor help")

	flag.StringVarP(&flags.MonitorIP, "monitorip", "i", "", "ip to be monitored")
	flag.StringVarP(&flags.MonitorService, "monitorservice", "s", "", "service name to be monitored")
	flag.Int32VarP(&flags.MonitorPid, "monitorpid", "p", -1, "pid to be monitored")
	flag.Int64VarP(&flags.MonitorInterval, "monitorinterval", "l", 1, "interval seconds to send monitor info")

	flag.StringVarP(&flags.DBUrl, "dburl", "d", "", "data base url")
	flag.StringVarP(&flags.Bucket, "bucket", "b", "", "data base bucket for influxdb")
	flag.StringVarP(&flags.Organization, "organization", "o", "", "data base org for influxdb")
	flag.StringVarP(&flags.Token, "token", "t", "", "data base token for influxdb")
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	flag.Parse()

	var monitorpid int32

	_, err := utils.CheckInK8s()
	if err == rest.ErrNotInCluster {
		//not running in k8s
		log.Info("running on bare metal")
		if flags.Help || flags.MonitorIP == "" || flags.MonitorPid == -1 ||
			flags.MonitorService == "" || flags.Bucket == "" || flags.DBUrl == "" ||
			flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			os.Exit(0)
		}

		monitorpid = flags.MonitorPid
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
	reporter = report.NewInfluxDBReporter(flags.DBUrl, flags.Organization, flags.Bucket, flags.Token, flags.MonitorInterval, processInfo)
	defer reporter.Close()
	reporter.Start(ctx)

	<-stopCh
	cancel()
	time.Sleep(5 * time.Second)
}
