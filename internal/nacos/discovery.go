package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

func CreateNacosClient() (naming_client.INamingClient, error) {
	clientConfig := getNacosClientConfig()
	serverConfigs := getNacosServerConfigs()

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getNacosClientConfig() constant.ClientConfig {
	return constant.ClientConfig{
		NamespaceId:         "ccac2447-2d93-4401-b0db-1555471bd09f",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            "nacos",
		Password:            "nacos",
	}
}

func getNacosServerConfigs() []constant.ServerConfig {
	return []constant.ServerConfig{
		{
			IpAddr: "43.198.104.122",
			Port:   8848,
			Scheme: "http",
		},
	}
}
