package router

import (
	"gate/api"
	"gate/auth"
	"github.com/gin-gonic/gin"
	"xmicro/internal/common/config"
)

// RegisterRouter 注册路由
func RegisterRouter() *gin.Engine {
	if config.Conf.ZapLog.Level == 1 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	//初始化grpc的client gate是做为grpc的客户端 去调用user grpc服务
	rpc.Init()
	r := gin.Default()
	r.Use(auth.Cors())
	userHandler := api.NewUserHandler()
	r.POST("/register", userHandler.Register)
	return r
}
