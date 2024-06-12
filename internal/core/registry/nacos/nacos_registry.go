package nacos

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
	"xmicro/internal/core/registry"
)

var once sync.Once
var reg registry.Registry

func NewNacosRegistry() registry.Registry {
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

func (r *Registry) Discover() ([]registry.InstanceInfo, error) {
	// Nacos服务发现逻辑
	return nil, nil
}

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

var (
	imServiceClients []im.ImServiceClient
	nextClientIndex  int32
	clientMu         sync.Mutex
)

// Setup initializes service discovery
func Setup() {
	client, err := newNamingClient()
	if err != nil {
		log.Printf("Nacos NewNamingClient error: %v", err)
		return
	}

	err = client.Subscribe(&vo.SubscribeParam{
		ServiceName: "grpc:im",
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				log.Printf("nacos subscribe service error: %v", err)
				return
			}

			updateGrpcClients(services)
		},
	})
	if err != nil {
		log.Printf("Nacos Subscribe error: %v", err)
		return
	}

	err = updateServiceInfo(client)
	if err != nil {
		log.Printf("Nacos failed to get initial service info: %v", err)
		return
	}
}

func updateGrpcClients(services []model.Instance) {
	clientMu.Lock()
	defer clientMu.Unlock()

	imServiceClients = nil
	for _, service := range services {
		log.Printf("ServiceName: %s, IP: %s, Port: %d, Metadata: %s\n",
			service.ServiceName, service.Ip, service.Port, service.Metadata)

		//conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "192.168.0.151", 9998), grpc.WithInsecure()) // 本地调试
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", service.Ip, service.Port), grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to gRPC service: %v", err)
			continue
		}
		imServiceClients = append(imServiceClients, im.NewImServiceClient(conn))
	}

	if len(imServiceClients) == 0 {
		log.Printf("No available gRPC services")
	}
}

func updateServiceInfo(client naming_client.INamingClient) error {
	services, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "grpc:im",
	})
	if err != nil {
		return fmt.Errorf("nacos GetService error: %w", err)
	}

	updateGrpcClients(services)

	return nil
}

func GetNextClient() im.ImServiceClient {
	clientMu.Lock()
	defer clientMu.Unlock()

	if len(imServiceClients) == 0 {
		//panic(exception.New(exception.ImNoAvailableGrpcClients, errors.New("")))
	}

	index := atomic.AddInt32(&nextClientIndex, 1)
	return imServiceClients[index%int32(len(imServiceClients))]
}
