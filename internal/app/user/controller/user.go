package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xmicro/internal/app/user/dto"
	"xmicro/internal/app/user/service"
	"xmicro/internal/constant"
	"xmicro/internal/x"
)

type UserController struct{}

func (that *UserController) Login(c *gin.Context) {
	req := dto.LoginReq{}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		x.HandleErr(c, x.ErrorModel(constant.ParamError))
		return
	}

	token, err := service.UserLogin(c, req)
	if err != nil {
		x.HandleErr(c, x.ErrorModel(constant.ServerError))
		return
	}

	x.Success(c, token)
}
