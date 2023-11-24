package bootstrap

import (
	"github.com/gin-gonic/gin"
	"xmicro/internal/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	smsCtl := new(controller.SmsController)

	// 发送百度智能云短信
	r.POST("/api/sms/send", smsCtl.Send)

	return r
}

//func RunServer() {
//	gin.SetMode(config.Conf.App.Mode)
//
//	r := SetupRouter()
//
//	logrus.Info("Server Start Success!")
//	if err := r.Run(":" + config.Conf.App.Port); err != nil {
//		logrus.Error("startup service failed, err:%v\n", err)
//	}
//}
