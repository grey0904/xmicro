package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net"
	"xmicro/internal/log"
)

var Client naming_client.INamingClient
var localIP string

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// Check the address type and if it is not a loopback, display it.
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("cannot find local IP address")
}

func RegistryToNacos() {
	var err error
	localIP, err = getLocalIP()
	if err != nil {
		log.Logger.Fatalf("failed to get local IP: %v", err)
	}

	Client, err = NewNamingClient()
	if err != nil {
		log.Logger.Fatalf("clients.CreateNamingClient error: %v", err)
	}

	// TODO 改为查询nacos中的配置
	serviceName := "grpc:user" // Set your service name here
	instance := vo.RegisterInstanceParam{
		Ip:          localIP, // Set your server IP here
		Port:        9998,    // Set your server port here
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

func DeregisterFromNacos() {
	instance := vo.DeregisterInstanceParam{
		Ip:          localIP,
		Port:        9998,
		ServiceName: "grpc:user",
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
