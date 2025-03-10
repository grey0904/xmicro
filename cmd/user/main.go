package main

import (
	"context"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"xmicro/internal/app/user"
	"xmicro/internal/common/config"
)

func main() {

	pflag.String("config", "", "path to config file (e.g., config/dev/config.yaml)")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	configPath := viper.GetString("config")
	if configPath == "" {
		log.Fatal("请提供配置文件路径，例如: --config config/dev/config.yaml")
		return
	}

	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	err = config.InitConfig("user", configPath)
	if err != nil {
		return
	}

	err = user.Run(context.Background())
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
