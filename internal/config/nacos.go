package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"net"
	"strconv"
	"xmicro/internal/app"
)

type Client struct {
	NacosConfig config_client.IConfigClient
}

func NewClient(c AppConfig) error {

	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         c.Nacos.NamespaceId, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           c.Nacos.TimeoutMs,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            c.Nacos.CacheDir,
		LogLevel:            "debug",
	}

	var serverConfigs []constant.ServerConfig

	// 至少一个ServerConfig
	for _, endpoint := range c.Nacos.Endpoints {
		ip, port, err := net.SplitHostPort(endpoint)
		if err != nil {
			fmt.Println("解析地址出错:", err)
		}

		value, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			fmt.Println("转换失败:", err)
		}

		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: ip,
			//ContextPath: "/nacos",
			Port:   value,
			Scheme: "http",
		})
	}

	// 创建动态配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	app.NacosClient = client

	return nil
}

type Db struct {
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
}
