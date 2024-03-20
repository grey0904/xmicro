package app

import (
	"errors"
	"flag"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"log"
	"xmicro/internal/config"
)

var (
	Rd     *redis.Client
	Db     *gorm.DB
	Nc     config_client.IConfigClient
	Config *config.AppConfig
)

func InitConfig() {
	InitMysqlConfig()
	InitRedisConfig()
}

func LoadConfig() {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configFile == "" {
		log.Fatalf("error:: %v", errors.New("123"))
	}

	viper.SetConfigFile(*configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error:: %v", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("error:: %v", err)
	}
}

func InitMysqlConfig() {

	content, err := Nc.GetConfig(vo.ConfigParam{
		DataId: "mysql.yaml",
	})
	if err != nil {
		log.Fatalf("initMysqlConfig NacosClient.GetConfig err: %v", err)
	}

	err = yaml.Unmarshal([]byte(content), &Config.Mysql)
	if err != nil {
		log.Fatalf("initMysqlConfig yaml.Unmarshal err: %v", err)
	}
}

func InitRedisConfig() {
	content, err := Nc.GetConfig(vo.ConfigParam{
		DataId: "redis.yaml",
	})
	if err != nil {
		log.Fatalf("initRedisConfig NacosClient.GetConfig err: %v", err)
	}

	err = yaml.Unmarshal([]byte(content), &Config.Redis)
	if err != nil {
		log.Fatalf("initRedisConfig yaml.Unmarshal err: %v", err)
	}
}
