package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func RegistryToNacos() error {
	clientConfig := getNacosClientConfig()
	serverConfigs := getNacosServerConfigs()

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return err
	}

	serviceName := "chat" // Set your service name here
	instance := vo.RegisterInstanceParam{
		Ip:          "127.0.0.1", // Set your server IP here
		Port:        9971,        // Set your server port here
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	success, err := client.RegisterInstance(instance)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to register gRPC service to Nacos")
	}

	fmt.Println("Registered gRPC service to Nacos successfully")
	return nil
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
