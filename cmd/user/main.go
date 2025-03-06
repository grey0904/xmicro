package main

import (
	"context"
	"log"
	"os"
	"xmicro/internal/app/user"
	"xmicro/internal/common/config/center"
)

func main() {

	err := center.InitConfig("user")
	if err != nil {
		return
	}
	
	// 启动 http、grpc 服务端
	err = user.RunV1(context.Background())
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
