package nacos

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strconv"
	"strings"
	"sync"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
	"xmicro/internal/core/registry"
)

var once sync.Once
var reg registry.Registry

func GetInstance() registry.Registry {
	once.Do(func() {
		client, err := NewNamingClient()
		if err != nil {
			logs.Error("NewNamingClient err:%v", err)
		}
		reg = &NacosRegistry{
			cli: client,
		}
	})
	return reg
}

type NacosRegistry struct {
	cli naming_client.INamingClient
}

func (r *NacosRegistry) Register(serviceName string) error {
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

func (r *NacosRegistry) Deregister(serviceName string) error {
	instance := vo.DeregisterInstanceParam{
		Ip:          config.Conf.Grpc.Host,
		Port:        config.Conf.Grpc.Port,
		ServiceName: serviceName,
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

func (n *NacosRegistry) Discover(serviceName string) ([]registry.InstanceInfo, error) {
	// Nacos服务发现逻辑
	return nil, nil
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
