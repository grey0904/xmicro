package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v2"
	"log"
	"xmicro/internal/nacos"
)

func main() {

	client, err := nacos.CreateConfigClient()
	if err != nil {
		log.Fatalf("Error CreateConfigClient: %v", err)
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "mysql.yaml",
	})
	if err != nil {
		log.Fatalf("Error client.GetConfig: %v", err)
	}

	// 解析YAML为Go结构
	var yamlConfig Db
	err = yaml.Unmarshal([]byte(content), &yamlConfig)
	if err != nil {
		log.Fatalf("Error yaml.Unmarshal: %v", err)
	}

	// 转换为JSON
	jsonData, err := json.Marshal(yamlConfig)
	if err != nil {
		log.Fatalf("Error json.Marshal: %v", err)
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
