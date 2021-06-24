package report

import (
	"context"
	"kube-gomonitor/agent/internal"
	"kube-gomonitor/agent/process"
	"kube-gomonitor/pkg"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

//InfluxDBReporter write info into InfluxDB
type InfluxDBReporter struct {
	Client         influxdb2.Client
	Writer         api.WriteAPI
	DBUrl          string
	Organization   string
	Bucket         string
	Token          string
	processInfo    *process.ProcessInfo
	reportInterval int64
}

func (reporter *InfluxDBReporter) Start(ctx context.Context) {
	go func() {
		log.Info("start reporter ...")
		for {
			select {
			case <-ctx.Done():
				log.Info("stop reporter ...")
				return
			default:
				reporter.processInfo.Update()
				reporter.Report()
				if reporter.processInfo.Fields.Status == "T" {
					return
				}
				time.Sleep(time.Duration(reporter.reportInterval) * time.Second)
			}
		}
	}()
}

func (reporter *InfluxDBReporter) Close() {
	reporter.Writer.Flush()
	reporter.Client.Close()
}

//Report send process info to db
func (reporter *InfluxDBReporter) Report() {
	// create point
	tags := pkg.Tags2Map(*reporter.processInfo.Tags)
	fileds := pkg.Fields2Map(*reporter.processInfo.Fields)
	p := influxdb2.NewPoint(
		"service-instance",
		tags,
		fileds,
		time.Now())

	reporter.Writer.WritePoint(p)
}

func NewInfluxDBReporter(flags *internal.CmdFlags, info *process.ProcessInfo) *InfluxDBReporter {
	client := influxdb2.NewClient(flags.DBUrl, flags.Token)
	writer := client.WriteAPI(flags.Organization, flags.Bucket)
	influxDbClient := InfluxDBReporter{
		client,
		writer,
		flags.DBUrl,
		flags.Organization,
		flags.Bucket,
		flags.Token,
		info,
		flags.MonitorInterval,
	}
	return &influxDbClient
}
