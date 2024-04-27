package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xmicro/internal/constant"
	"xmicro/internal/service/user/dto"
	"xmicro/internal/service/user/service"
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

func (that *UserController) Orders(c *gin.Context) {
	req := dto.UserOrdersReq{}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		x.HandleErr(c, x.ErrorModel(constant.ParamError))
		return
	}

	token, err := service.UserOrders(c, req)
	if err != nil {
		x.HandleErr(c, x.ErrorModel(constant.ServerError))
		return
	}

	x.Success(c, token)
}
