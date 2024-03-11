package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

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

func DeregisterFromNacos() error {
	clientConfig := getNacosClientConfig()
	serverConfigs := getNacosServerConfigs()

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return err
	}

	serviceName := "chat"
	instance := vo.DeregisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        9971,
		ServiceName: serviceName,
	}
	success, err := client.DeregisterInstance(instance)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to register gRPC service to Nacos")
	}

	fmt.Println("Unregistered gRPC service from Nacos successfully")
	return nil
}

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
