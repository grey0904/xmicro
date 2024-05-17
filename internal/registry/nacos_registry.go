package registry

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strconv"
	"strings"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
)

type NacosRegistry struct {
	cli naming_client.INamingClient
}

func (r *NacosRegistry) Register(appName string) error {
	serviceName := "grpc:" + appName // Set your service name here
	instance := vo.RegisterInstanceParam{
		Ip:          config.Conf.Grpc.Host, // Set your server IP here
		Port:        config.Conf.Grpc.Port, // Set your server port here
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	success, err := r.cli.RegisterInstance(instance)
	if err != nil {
		return fmt.Errorf("clients.RegisterInstance error: %v", err)
	}
	if !success {
		return fmt.Errorf("clients.RegisterInstance was not successful")
	}

	logs.Info("Registered gRPC service to Nacos successfully")
	return nil
}

func (r *NacosRegistry) Deregister(appName string) error {
	instance := vo.DeregisterInstanceParam{
		Ip:          config.Conf.Grpc.Host,
		Port:        config.Conf.Grpc.Port,
		ServiceName: "grpc:" + appName,
	}
	success, err := r.cli.DeregisterInstance(instance)
	if err != nil {
		return fmt.Errorf("failed to deregister from Nacos: %v", err)
	}
	if !success {
		return fmt.Errorf("failed to deregister from Nacos: unsuccess")
	}

	logs.Info("Successfully deregistered from Nacos")
	return nil
}

func (r *NacosRegistry) GetService(appName string) error {
	appName = "grpc:" + appName
	client, err := NewNamingClient()
	if err != nil {
		panic(err)
	}

	err = client.Subscribe(&vo.SubscribeParam{
		ServiceName: appName,
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				panic(err)
			}

			updateGrpcClients(services)
		},
	})
	if err != nil {
		panic(err)
	}

	services, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: appName,
	})
	if err != nil {
		panic(err)
	}

	updateGrpcClients(services)
	return nil
}

// ========= 工厂模块  =========
type NacosFactory struct {
}

func (fac *NacosFactory) CreateRegistrar() (Registrar, error) {
	var reg Registrar
	var err error
	clientOnce.Do(func() {
		Client, err = NewNamingClient()
	})
	if err != nil {
		return nil, fmt.Errorf("clients.CreateNamingClient error: %v", err)
	}

	reg = &NacosRegistry{
		cli: Client,
	}
	return reg, nil
}

func NewNamingClient() (naming_client.INamingClient, error) {
	var (
		sc = make([]constant.ServerConfig, 0)
		nc = config.LocalConf.Nacos
	)

	for _, value := range nc.Endpoints {
		vs := strings.Split(value, ":")

		if len(vs) < 2 {
			return nil, errors.New("endpoints configuration error")
		}

		port, err := strconv.ParseUint(vs[1], 10, 64)
		if err != nil {
			return nil, errors.New("endpoints configuration error")
		}

		sc = append(sc, constant.ServerConfig{
			IpAddr: vs[0],
			Port:   port,
		})
	}

	cc := constant.ClientConfig{
		NamespaceId:         nc.NamespaceId,
		TimeoutMs:           nc.TimeoutMs,
		NotLoadCacheAtStart: true,
		CacheDir:            nc.CacheDir,
		Username:            nc.Username,
		Password:            nc.Password,
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Nacos client: %w", err)
	}

	return client, nil
}
