package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SmsController struct{}

func (that *SmsController) Send(c *gin.Context) {
	mobile := c.PostForm("mobile")

	logrus.Infof("%s短信验证码发送成功", mobile)
}
