package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"xmicro/internal/app"
	"xmicro/internal/common/config"
	"xmicro/internal/common/metrics"
	"xmicro/internal/nacos"
	"xmicro/internal/service/order/router"
)

func main() {
	// 1.加载配置
	config.InitConfig() // 用viper将config/dev/nacos-local.yaml文件的数据解析到 AppConfig 结构体

	//2.启动监控
	go func() {
		err := metrics.Serve(fmt.Sprintf("0.0.0.0:%d", config.Conf.MetricPort))
		if err != nil {
			panic(err)
		}
	}()

	//3.启动 http、grpc 服务端
	err := app.Run(context.Background())
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	// 注册服务
	nacos.RegistryToNacos("order")

	// 启动本地服务
	router.RunServer()

	// 执行取消注册操作
	nacos.DeregisterFromNacos("order")

	// 等待一段时间确保异步处理完成
	time.Sleep(2 * time.Second)
}
