package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"xmicro/internal/app/user/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userCtl := new(controller.UserController)

	r.POST("/login", userCtl.Login)

	return r
}

func RunServer() {
	gin.SetMode(gin.DebugMode)

	r := SetupRouter()

	logrus.Info("Server Start Success!")

	if err := r.Run(":8080"); err != nil {
		logrus.Error("startup service failed, err:%v\n", err)
	}
}
