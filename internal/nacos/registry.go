package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net"
)

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

func RegistryToNacos() error {
	localIP, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP: %v", err)
	}

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
		Ip:          localIP, // Set your server IP here
		Port:        9998,    // Set your server port here
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
	localIP, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP: %v", err)
	}

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
		Ip:          localIP,
		Port:        9998,
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
