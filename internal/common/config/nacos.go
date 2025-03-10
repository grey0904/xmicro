package config

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosCenter struct {
	client config_client.IConfigClient
}

func NewNacosCenter(cfg interface{}) (*NacosCenter, error) {
	nacosConfig, ok := cfg.(NacosClientConfig)
	if !ok {
		return nil, fmt.Errorf("invalid nacos config type")
	}

	sc := make([]constant.ServerConfig, 0)
	for _, endpoint := range nacosConfig.Endpoints {
		host, port, err := parseEndpoint(endpoint)
		if err != nil {
			return nil, err
		}
		sc = append(sc, constant.ServerConfig{
			IpAddr: host,
			Port:   port,
		})
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.NamespaceId,
		Username:            nacosConfig.Username,
		Password:            nacosConfig.Password,
		TimeoutMs:           nacosConfig.TimeoutMs,
		NotLoadCacheAtStart: nacosConfig.NotLoadCacheAtStart,
		LogDir:              nacosConfig.LogDir,
		CacheDir:            nacosConfig.CacheDir,
		LogLevel:            nacosConfig.LogLevel,
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return nil, fmt.Errorf("create nacos client failed: %v", err)
	}

	return &NacosCenter{client: client}, nil
}

func (n *NacosCenter) GetConfig(key string) (string, error) {
	content, err := n.client.GetConfig(vo.ConfigParam{
		DataId: key,
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		return "", fmt.Errorf("get config from nacos failed: %v", err)
	}
	return content, nil
}

func (n *NacosCenter) WatchConfig(key string, onChange func(string)) error {
	return n.client.ListenConfig(vo.ConfigParam{
		DataId: key,
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			onChange(data)
		},
	})
}

func (n *NacosCenter) Close() error {
	// Nacos client doesn't provide Close method
	return nil
}

// parseEndpoint 解析 endpoint 字符串为 host 和 port
func parseEndpoint(endpoint string) (string, uint64, error) {
	var host string
	var port uint64
	_, err := fmt.Sscanf(endpoint, "%s:%d", &host, &port)
	if err != nil {
		return "", 0, fmt.Errorf("invalid endpoint format: %s", endpoint)
	}
	return host, port, nil
}
