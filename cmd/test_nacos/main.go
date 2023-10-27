package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v2"
)

func main() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "0f9a9129-7fcb-40f4-9138-5d3720c5bd4b", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			//ContextPath: "/nacos",
			Port:   8848,
			Scheme: "http",
		},
	}

	// 创建动态配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "mysql.yaml",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// 解析YAML为Go结构
	var yamlConfig Db
	err = yaml.Unmarshal([]byte(content), &yamlConfig)
	if err != nil {
		fmt.Println(err)
	}

	// 转换为JSON
	jsonData, err := json.Marshal(yamlConfig)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s", jsonData)
}

type Db struct {
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
}
