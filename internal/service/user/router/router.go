package router

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"xmicro/internal/nacos"
	"xmicro/internal/service/user/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("user/")
	userCtl := new(controller.UserController)

	api.POST("/login", userCtl.Login)
	api.POST("/orders", userCtl.Orders)

	// 创建一个GET路由，用于获取服务实例信息
	api.GET("/instance", func(c *gin.Context) {
		instances, err := nacos.Client.SelectAllInstances(vo.SelectAllInstancesParam{
			ServiceName: "grpc:user",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取服务实例信息"})
			return
		}
		c.JSON(http.StatusOK, instances)
	})

	return r
}

func RunServer() {
	var err error
	gin.SetMode(gin.DebugMode)

	r := SetupRouter()

	logrus.Info("Server Start Success!")

	// 启动GIN服务器
	go func() {
		err = r.Run(":8080")
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
