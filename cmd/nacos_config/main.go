package main

import (
	"xmicro/internal/app"
	"xmicro/internal/database"
)

func main() {
	app.InitConfig()
	database.InitRedis()
	database.InitMysql()
}
