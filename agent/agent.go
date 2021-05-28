package main

import (
	"context"
	"gomonitor/agent/packet"
	"gomonitor/agent/process"
	"gomonitor/agent/report"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

type cmdFlags struct {
	Help            bool
	MonitorIP       string
	MonitorPid      int32
	MonitorService  string
	MonitorInterval int64
}

var flags cmdFlags

func init() {
	flag.StringVarP(&flags.MonitorIP, "monitorip", "i", "", "ip to be monitored")
	flag.StringVarP(&flags.MonitorService, "monitorservice", "s", "", "service name to be monitored")
	flag.Int32VarP(&flags.MonitorPid, "monitorpid", "p", -1, "pid to be monitored")
	flag.Int64VarP(&flags.MonitorInterval, "monitorinterval", "l", 1, "interval seconds to send monitor info")
	flag.BoolVarP(&flags.Help, "help", "h", false, "gomonitor help")
}

func main() {
	flag.Parse()

	if flags.Help || flags.MonitorIP == "" || flags.MonitorPid == -1 || flags.MonitorService == "" {
		flag.Usage()
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go packet.NetSniff(ctx, flags.MonitorIP)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				processInfo, err := process.ProcessStat(flags.MonitorPid, flags.MonitorService)
				if err != nil {
					log.Error(err)
					continue
				}
				processInfo.IP = flags.MonitorIP
				report.Report(processInfo)
				time.Sleep(time.Duration(flags.MonitorInterval) * time.Second)
			}
		}
	}()

	<-stopCh
	cancel()
	time.Sleep(5 * time.Second)
}
