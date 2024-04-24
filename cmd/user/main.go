package main

import (
	"time"
	"xmicro/internal/app"
	"xmicro/internal/database"
	"xmicro/internal/log"
	"xmicro/internal/nacos"
	"xmicro/internal/services/user/router"
)

func main() {
	// 加载配置
	app.LoadConfig()        // 用viper将config/dev/nacos-local.yaml文件的数据解析到 AppConfig 结构体
	nacos.NewConfigClient() // 用 AppConfig 中的Nacos配置信息创建“配置中心客户端”
	app.InitConfig()        // 从Nacos上获取Mysql、Redis等配置，并解析给对应的 AppConfig 里面的结构体

	// 初始化组件
	database.InitRedis()
	database.InitMysql()
	log.InitLogger()

	//nacos.RegistryToNacos("user")
	//rpc.UserRpcRegister()
	nacos.DiscoveryFromNacos("order")

	// goroutine 启动本地服务
	router.RunServer()

	// 执行取消注册操作
	//nacos.DeregisterFromNacos("user")

	// 等待一段时间确保异步处理完成
	time.Sleep(2 * time.Second)
}
