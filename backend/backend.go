package main

import (
	"gomonitor/backend/server"
	"gomonitor/backend/server/bare"
	"gomonitor/backend/server/k8s"
	"gomonitor/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/rest"
)

var (
	flags server.ServerFlags
	//nacos默认配置
	logDir     string = "/tmp/nacos/log"
	cacheDir   string = "/tmp/nacos/cache"
	rotateTime string = "12h"
	maxAge     int64  = 3
	logLevel   string = "info"
	timeoutMs  uint64 = 5000
	//server
	monitorServer server.MonitorServer
)

func init() {
	flag.Uint64VarP(&flags.NacosPort, "nacosport", "p", 8848, "nacos server port")
	flag.StringVarP(&flags.NacosIP, "nacosip", "i", "", "nacos server ip")

	flag.StringVarP(&flags.MonitorServiceGroup, "group", "g", "DEFAULT_GROUP", "monitor service group")
	flag.StringSliceVarP(&flags.MonitorServices, "monitorservices", "s", nil, "monitor service names")
	flag.StringVarP(&flags.NamespaceId, "namespace", "n", "public", "nacos namespace id (not namespace name)")
	flag.Uint64VarP(&flags.Interval, "interval", "l", 3, "agent monitor interval")

	flag.StringVarP(&flags.DBUrl, "dburl", "d", "", "data base url")
	flag.StringVarP(&flags.Bucket, "bucket", "b", "", "data base bucket for influxdb")
	flag.StringVarP(&flags.Organization, "organization", "o", "", "data base org for influxdb")
	flag.StringVarP(&flags.Token, "token", "t", "", "data base token for influxdb")
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	_, err := utils.CheckInK8s()

	if err == rest.ErrNotInCluster {
		flag.Parse()
		if flags.NacosIP == "" || flags.MonitorServices == nil || flags.DBUrl == "" ||
			flags.Bucket == "" || flags.Organization == "" || flags.Token == "" {
			flag.Usage()
			return
		}

		log.Info("running on bare metal")
		//get nacos config
		sc := []constant.ServerConfig{
			{
				IpAddr: flags.NacosIP,
				Port:   flags.NacosPort,
			},
		}
		cc := constant.ClientConfig{
			NamespaceId:         flags.NamespaceId, //namespace id
			TimeoutMs:           timeoutMs,
			NotLoadCacheAtStart: true,
			LogDir:              logDir,
			CacheDir:            cacheDir,
			RotateTime:          rotateTime,
			MaxAge:              maxAge,
			LogLevel:            logLevel,
		}

		monitorServer, err = bare.NewBareMetalServer(vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		}, &flags)

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
		monitorServer = k8s.NewKServer(flags.DBUrl, flags.Bucket, flags.Organization, flags.Token)
		monitorServer.Start()

		<-signalChan
		time.Sleep(5 * time.Second)
	}
}
