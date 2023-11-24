package config

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"xmicro/internal/log"
)

type Nacos struct {
	Endpoints   []string `yaml:"endpoints"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	TimeoutMs   uint64   `yaml:"timeoutMs"`
	NamespaceId string   `yaml:"namespaceId"`
	CacheDir    string   `yaml:"cacheDir"`
}

type AppConfig struct {
	Nacos Nacos `yaml:"nacos"`
}

func LoadConfig() error {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configFile == "" {
		log.Logger.Error("error:", errors.New("123"))
	}

	viper.SetConfigFile(*configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Logger.Error("error:", err)
		return err
	}

	serverHost := viper.GetString("nacos.username")
	fmt.Println(serverHost)

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		log.Logger.Error("error:", err)
		return err
	}

	err = NewClient(appConfig)
	if err != nil {
		log.Logger.Error("error:", err)
	}

	return nil
}
