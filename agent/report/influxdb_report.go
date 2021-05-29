package report

import (
	"gomonitor/agent/process"
	"gomonitor/agent/utils"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

//InfluxDBReporter write info into InfluxDB
type InfluxDBReporter struct {
	Client       influxdb2.Client
	Writer       api.WriteAPI
	DBUrl        string
	Organization string
	Bucket       string
	Token        string
}

func (reporter *InfluxDBReporter) Close() {
	reporter.Writer.Flush()
	reporter.Client.Close()
}

//Report send process info to db
func (reporter *InfluxDBReporter) Report(processInfo *process.ProcessInfo) {
	// create point
	tags := utils.Tags2Map(*processInfo.Tags)
	fileds := utils.Fields2Map(*processInfo.Fields)
	p := influxdb2.NewPoint(
		"service-instance",
		tags,
		fileds,
		time.Now())

	reporter.Writer.WritePoint(p)
}

func NewInfluxDBReporter(url string, organization string,
	bucket string, token string) *InfluxDBReporter {
	client := influxdb2.NewClient(url, token)
	writer := client.WriteAPI(organization, bucket)
	influxDbClient := InfluxDBReporter{
		client,
		writer,
		url,
		organization,
		bucket,
		token,
	}
	return &influxDbClient
}
