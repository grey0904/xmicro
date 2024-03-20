package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"xmicro/internal/app"
)

func InitMysql() {
	c := app.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("initMysqlConfig yaml.Unmarshal err: %v", err)
	}
	app.Db = db
}
