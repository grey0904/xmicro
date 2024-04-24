package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"xmicro/internal/log"
	"xmicro/internal/utils/u_ip"
)

var (
	Client  naming_client.INamingClient
	localIP string
)

func RegistryToNacos(svcName string) {
	var (
		err error
	)

	localIP, err = u_ip.GetLocalIP()
	if err != nil {
		log.Logger.Fatalf("failed to get local IP: %v", err)
	}

	Client, err = NewNamingClient()
	if err != nil {
		log.Logger.Fatalf("clients.CreateNamingClient error: %v", err)
	}

	// TODO 改为查询nacos中的配置
	serviceName := "grpc:" + svcName // Set your service name here
	instance := vo.RegisterInstanceParam{
		Ip:          localIP, // Set your server IP here
		Port:        9997,    // Set your server port here
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	success, err := Client.RegisterInstance(instance)
	if err != nil {
		log.Logger.Fatalf("clients.RegisterInstance error: %v", err)
	}
	if !success {
		log.Logger.Fatalf("clients.RegisterInstance error: %v", err)
	}

	fmt.Println("Registered gRPC service to Nacos successfully")
}

func DeregisterFromNacos(svcName string) {
	instance := vo.DeregisterInstanceParam{
		Ip:          localIP,
		Port:        9997,
		ServiceName: "grpc:" + svcName,
	}
	success, err := Client.DeregisterInstance(instance)
	if err != nil {
		log.Logger.Fatalf("Failed to deregister from Nacos: %v", err.Error)
	}
	if !success {
		log.Logger.Fatalf("Failed to deregister from Nacos: unsuccess")
	}

	log.Logger.Info("Successfully deregistered from Nacos")
}
