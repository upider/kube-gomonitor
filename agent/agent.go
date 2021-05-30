package main

import (
	"context"
	"gomonitor/agent/packet"
	"gomonitor/agent/process"
	"gomonitor/agent/report"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

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
	flag.StringVarP(&flags.MonitorIP, "monitorip", "i", "", "ip to be monitored")
	flag.StringVarP(&flags.MonitorService, "monitorservice", "s", "", "service name to be monitored")
	flag.Int32VarP(&flags.MonitorPid, "monitorpid", "p", -1, "pid to be monitored")
	flag.Int64VarP(&flags.MonitorInterval, "monitorinterval", "l", 1, "interval seconds to send monitor info")
	flag.BoolVarP(&flags.Help, "help", "h", false, "gomonitor help")

	flag.StringVarP(&flags.DBUrl, "dburl", "d", "", "ip to be monitored")
	flag.StringVarP(&flags.Bucket, "bucket", "b", "", "ip to be monitored")
	flag.StringVarP(&flags.Organization, "organization", "o", "", "ip to be monitored")
	flag.StringVarP(&flags.Token, "token", "t", "", "ip to be monitored")
}

func main() {
	flag.Parse()

	if flags.Help || flags.MonitorIP == "" || flags.MonitorPid == -1 ||
		flags.MonitorService == "" || flags.Bucket == "" || flags.DBUrl == "" ||
		flags.Organization == "" || flags.Token == "" {
		flag.Usage()
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go packet.NetSniff(ctx, flags.MonitorIP)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	processInfo, err := process.NewProcessInfo(flags.MonitorPid, flags.MonitorService, flags.MonitorIP)
	if err != nil {
		log.Error(err)
		return
	}

	reporter = report.NewInfluxDBReporter(flags.DBUrl, flags.Organization, flags.Bucket, flags.Token, flags.MonitorInterval, processInfo)
	defer reporter.Close()

	reporter.Start(ctx)

	<-stopCh
	time.Sleep(5 * time.Second)
}
