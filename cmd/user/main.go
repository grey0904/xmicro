package main

import (
	"xmicro/internal/app/user/router"
	"xmicro/internal/config"
	"xmicro/internal/database"
	"xmicro/internal/log"
)

func main() {

	log.InitLogger()

	if err := config.LoadConfig(); err != nil {
		log.Logger.Error("Load configs json error:", err)
		return
	}

	err := database.InitRedis()
	if err != nil {
		return
	}

	//err = bootstrap.InitMysql()
	//if err != nil {
	//	return
	//}

	router.RunServer()
}
