package main

import (
	"xmicro/internal/app"
	"xmicro/internal/app/user/router"
	"xmicro/internal/database"
	"xmicro/internal/log"
	"xmicro/internal/nacos"
)

func main() {
	app.LoadConfig()
	nacos.CreateConfigClient()
	app.InitConfig()
	database.InitRedis()
	database.InitMysql()
	log.InitLogger()

	router.RunServer()
}
