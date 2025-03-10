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
	//Discover() ([]InstanceInfo, error)
}

// InitRegistry 根据配置文件初始化服务注册实现
func InitRegistry() (Registry, error) {
	var reg Registry
	kind := config.LocalConf.Type
	switch kind {
	case "nacos":
		reg = nacos.NewNacosRegistry()
	case "etcd":
		reg = etcd.NewEtcdRegistry()
	default:
		return nil, fmt.Errorf("unknown registry kind: %s", kind)
	}
	return reg, nil
}
