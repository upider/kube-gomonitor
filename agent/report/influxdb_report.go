package report

import (
	"context"
	"gomonitor/agent/process"
	"gomonitor/agent/utils"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
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
		for {
			select {
			case <-ctx.Done():
				return
			default:
				reporter.processInfo.Update()
				reporter.Report()
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
	tags := utils.Tags2Map(*reporter.processInfo.Tags)
	fileds := utils.Fields2Map(*reporter.processInfo.Fields)
	p := influxdb2.NewPoint(
		"service-instance",
		tags,
		fileds,
		time.Now())

	reporter.Writer.WritePoint(p)
}

func NewInfluxDBReporter(url string, organization string,
	bucket string, token string, interval int64, info *process.ProcessInfo) *InfluxDBReporter {
	client := influxdb2.NewClient(url, token)
	writer := client.WriteAPI(organization, bucket)
	influxDbClient := InfluxDBReporter{
		client,
		writer,
		url,
		organization,
		bucket,
		token,
		info,
		interval,
	}
	return &influxDbClient
}
