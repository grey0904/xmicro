package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"xmicro/internal/app/user"
	"xmicro/internal/common/config"
	"xmicro/internal/common/metrics"
)

func main() {

	// 1.加载配置
	config.InitConfig("user") // 用viper将config/dev/nacos-local.yaml文件的数据解析到 AppConfig 结构体

	//2.启动监控
	go func() {
		err := metrics.Serve(fmt.Sprintf("0.0.0.0:%d", config.Conf.MetricPort))
		if err != nil {
			panic(err)
		}
	}()

	//3.启动 http、grpc 服务端
	err := user.RunV2(context.Background())
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
