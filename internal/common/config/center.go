package config

import (
	"fmt"
)

// ConfigCenter 配置中心接口
type ConfigCenter interface {
	// GetConfig 获取配置
	GetConfig(key string) (string, error)
	// WatchConfig 监听配置变化
	WatchConfig(key string, onChange func(string)) error
	// Close 关闭连接
	Close() error
}

// NewConfigCenter 创建配置中心客户端
func NewConfigCenter(centerType string, config interface{}) (ConfigCenter, error) {
	switch centerType {
	case "nacos":
		return NewNacosCenter(config)
	case "etcd":
		return NewEtcdCenter(config)
	default:
		return nil, fmt.Errorf("unsupported config center type: %s", centerType)
	}
}
