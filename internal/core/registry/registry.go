package registry

import (
	"fmt"
	"xmicro/internal/common/config"
	"xmicro/internal/core/registry/etcd"
	"xmicro/internal/core/registry/nacos"
)

type InstanceInfo struct {
	ID       string
	Address  string
	Port     int
	Metadata map[string]string
}

type Registry interface {
	Register() error
	Deregister() error
	Discover() ([]InstanceInfo, error)
}

// InitRegistry 根据配置文件初始化服务注册实现
func InitRegistry(conf config.RegistryConfig) (Registry, error) {
	var reg Registry
	config.Conf.Services
	switch conf.Type {
	case "nacos":
		reg = nacos.NewNacosRegistry(conf.Nacos)
	case "etcd":
		reg = etcd.NewEtcdRegistry(conf.Etcd)
	default:
		return nil, fmt.Errorf("unknown registry type: %s", conf.Type)
	}
	return reg, nil
}
