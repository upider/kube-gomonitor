package main

import (
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func bareRegistry() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "192.168.37.204",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	pid := strconv.Itoa(os.Getpid())
	client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "192.168.37.204",
		Port:        8848,
		ServiceName: "test-service",
		GroupName:   "DEFAULT_GROUP",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"pid": pid},
	})
}
