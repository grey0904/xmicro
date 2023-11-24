package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	// 设置配置文件的名称（不包含扩展名）
	viper.SetConfigName("config")
	// 设置配置文件的类型
	viper.SetConfigType("yaml")
	// 添加配置文件所在的路径（可选）
	viper.AddConfigPath("./config/dev")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("无法读取配置文件:", err)
		return
	}

	// 读取配置项的值
	serverHost := viper.GetString("server.host")
	serverPort := viper.GetInt("server.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")

	// 打印配置项的值
	fmt.Printf("Server Host: %s\n", serverHost)
	fmt.Printf("Server Port: %d\n", serverPort)
	fmt.Printf("Database User: %s\n", dbUser)
	fmt.Printf("Database Password: %s\n", dbPassword)
}
