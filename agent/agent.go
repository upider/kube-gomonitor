package main

import (
	"context"
	"encoding/json"
	"gomonitor/agent/packet"
	"gomonitor/agent/process"
	"gomonitor/agent/report"
	"gomonitor/utils"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"

	"syscall"

	flag "github.com/spf13/pflag"
)

type monitorConfig struct {
	Pid         int32  `json:"pid"`
	ServiceName string `json:"serviceName"`
	IP          string `json:"ip"`
}

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
	flags       cmdFlags
	reporter    report.Reporter
	configPath  string = "/tmp/monitor-config.json"
	processInfo *process.ProcessInfo
	content     []byte
	configs     monitorConfig
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

	_, err := utils.CheckInK8s()
	if err == rest.ErrNotInCluster {
		//not running in k8s
		flag.Parse()

		if flags.Help || flags.MonitorIP == "" || flags.MonitorPid == -1 ||
			flags.MonitorService == "" || flags.Bucket == "" || flags.DBUrl == "" ||
			flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			os.Exit(0)
		}

		go packet.NetSniff(ctx, flags.MonitorIP)

		processInfo, err = process.NewProcessInfo(flags.MonitorPid, flags.MonitorService, flags.MonitorIP)
		if err != nil {
			log.Error(err)
			return
		}
		reporter = report.NewInfluxDBReporter(flags.DBUrl, flags.Organization, flags.Bucket, flags.Token, flags.MonitorInterval, processInfo)

	} else {
		//running in k8s
		if flags.Help || flags.Bucket == "" || flags.DBUrl == "" ||
			flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			os.Exit(0)
		}

		//读取/tmp/monitor-config.json
		fin, err := os.Open(configPath)
		if err != nil {
			panic(err)
		}
		content, err = ioutil.ReadAll(fin)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(content, &configs)
		if err != nil {
			panic(err)
		}

		processInfo, err = process.NewProcessInfo(configs.Pid, configs.ServiceName, configs.IP)
		if err != nil {
			log.Error(err)
			return
		}
		reporter = report.NewInfluxDBReporter(flags.DBUrl, flags.Organization, flags.Bucket, flags.Token, flags.MonitorInterval, processInfo)

	}

	defer reporter.Close()

	reporter.Start(ctx)

	<-stopCh
	cancel()
	time.Sleep(5 * time.Second)
}
