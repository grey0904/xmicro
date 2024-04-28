package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
)

var (
	Client naming_client.INamingClient
)

func RegistryToNacos(appName string) {
	var (
		err error
	)

	Client, err = NewNamingClient()
	if err != nil {
		logs.Fatal("clients.CreateNamingClient error: %v", err)
	}

	// TODO 改为查询nacos中的配置
	serviceName := "grpc:" + appName // Set your service name here
	instance := vo.RegisterInstanceParam{
		Ip:          config.Conf.Grpc.Host, // Set your server IP here
		Port:        config.Conf.Grpc.Port, // Set your server port here
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	success, err := Client.RegisterInstance(instance)
	if err != nil {
		logs.Fatal("clients.RegisterInstance error: %v", err)
	}
	if !success {
		logs.Fatal("clients.RegisterInstance error: %v", err)
	}

	fmt.Println("Registered gRPC service to Nacos successfully")
}

func DeregisterFromNacos(appName string) {
	instance := vo.DeregisterInstanceParam{
		Ip:          config.Conf.Grpc.Host,
		Port:        config.Conf.Grpc.Port,
		ServiceName: "grpc:" + appName,
	}
	success, err := Client.DeregisterInstance(instance)
	if err != nil {
		logs.Fatal("Failed to deregister from Nacos: %v", err.Error)
	}
	if !success {
		logs.Fatal("Failed to deregister from Nacos: unsuccess")
	}

	logs.Fatal("Successfully deregistered from Nacos")
}
