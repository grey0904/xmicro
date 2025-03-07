package router

import (
	"github.com/gin-gonic/gin"
	"xmicro/internal/app/gate/api"
	"xmicro/internal/app/gate/auth"
	"xmicro/internal/common/config/center"
	"xmicro/internal/common/rpc/discovery"
)

// RegisterRouter 注册路由
func RegisterRouter() *gin.Engine {
	if config.Conf.ZapLog.Level == -1 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	//初始化grpc的client gate是做为grpc的客户端 去调用user grpc服务
	discovery.Init()
	r := gin.Default()
	r.Use(auth.Cors())
	userHandler := api.NewUserHandler()
	r.POST("/register", userHandler.Register)
	return r
}
