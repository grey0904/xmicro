package router

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"xmicro/internal/services/order/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	orderCtl := new(controller.OrderController)

	r.POST("order/list", orderCtl.List)

	return r
}

func RunServer() {
	var err error
	gin.SetMode(gin.DebugMode)

	r := SetupRouter()

	logrus.Info("Server Start Success!")

	// 启动GIN服务器
	go func() {
		err = r.Run(":8081")
		if err != nil {
			logrus.Error("启动GIN服务器失败: " + err.Error())
			return
		}
	}()

	// 等待终止信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
