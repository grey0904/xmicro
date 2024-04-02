package main

import (
	"context"
	"google.golang.org/grpc"
	"time"
	"xmicro/internal/app"
	"xmicro/internal/app/user/router"
	"xmicro/internal/database"
	"xmicro/internal/log"
	"xmicro/internal/nacos"
	"xmicro/internal/pb/user"
)

type userServer struct {
	user.UnimplementedUserServiceServer
}

func (s *userServer) GetUserInfo(ctx context.Context, in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	return &user.GetUserInfoResponse{
		UserId:   "1",
		UserName: "ben",
		Email:    "benben@gmail.com",
		Age:      "22",
	}, nil
}

func main() {
	// 加载配置
	app.LoadConfig()        // 用viper将config/dev/nacos-local.yaml文件的数据解析到 AppConfig 结构体
	nacos.NewConfigClient() // 用 AppConfig 中的Nacos配置信息创建“配置中心客户端”
	app.InitConfig()        // 从Nacos上获取Mysql、Redis等配置，并解析给对应的 AppConfig 里面的结构体

	// 初始化组件
	database.InitRedis()
	database.InitMysql()
	log.InitLogger()

	// 注册服务
	nacos.RegistryToNacos()
	user.RegisterUserServiceServer(grpc.NewServer(), &userServer{})

	// goroutine 启动本地服务
	router.RunServer()

	// 执行取消注册操作
	nacos.DeregisterFromNacos()

	// 等待一段时间确保异步处理完成
	time.Sleep(2 * time.Second)
}
