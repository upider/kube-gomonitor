package pkg

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// GetNacosNamingClient 获取nacos NamingClients
func GetNacosNamingClients(nacosIPs []string, nacosPort uint64, namespaces []string) ([]naming_client.INamingClient, error) {
	var nacosNamingClients []naming_client.INamingClient

	var serverConfigs []constant.ServerConfig
	for _, ip := range nacosIPs {
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(
			ip,
			nacosPort,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		))
	}

	clientConfig := *constant.NewClientConfig(
		// constant.WithNamespaceId(namespace), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel("info"),
	)

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		return nil, err
	}

	nacosNamingClients = append(nacosNamingClients, namingClient)

	return nacosNamingClients, nil
}
