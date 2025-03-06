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
	"xmicro/internal/common/config/center"
	"xmicro/internal/common/logs"
)

var once sync.Once
var reg *Registry

func NewNacosRegistry() *Registry {
	once.Do(func() {
		client, err := newNamingClient()
		if err != nil {
			logs.Error("NewNamingClient err:%v", err)
			return
		}
		reg = &Registry{
			cli: client,
		}
	})
	return reg
}

type Registry struct {
	cli naming_client.INamingClient
}

func (r *Registry) Register() error {
	instance := vo.RegisterInstanceParam{
		Ip:          config.Conf.ServerRpc.Host, // Set your server IP here
		Port:        config.Conf.ServerRpc.Port, // Set your server port here
		ServiceName: config.LocalConf.AppName,
		Weight:      config.Conf.Nacos.Weight,
		Enable:      true,
		Healthy:     true,
		//Ephemeral:   true, // 是否是临时实例，true为是
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

func (r *Registry) Deregister() error {
	instance := vo.DeregisterInstanceParam{
		Ip:          config.Conf.ServerRpc.Host,
		Port:        config.Conf.ServerRpc.Port,
		ServiceName: config.LocalConf.AppName,
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

//func (r *Registry) Discover() ([]registry.InstanceInfo, error) {
//	// Nacos服务发现逻辑
//	return nil, nil
//}

func newNamingClient() (naming_client.INamingClient, error) {
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
