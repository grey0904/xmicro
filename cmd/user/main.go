package main

import (
	"xmicro/bootstrap"
	"xmicro/internal/config"
	"xmicro/internal/log"
)

func main() {

	log.InitLogger()

	if err := config.LoadConfig(); err != nil {
		log.Logger.Error("Load configs json error:", err)
		return
	}

	err := bootstrap.InitRedis()
	if err != nil {
		return
	}

	//err = bootstrap.InitMysql()
	//if err != nil {
	//	return
	//}

	//bootstrap.RunServer()
}
