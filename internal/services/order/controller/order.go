package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xmicro/internal/constant"
	"xmicro/internal/services/order/dto"
	"xmicro/internal/services/order/service"
	"xmicro/internal/x"
)

type OrderController struct{}

func (that *OrderController) List(c *gin.Context) {
	req := dto.OrderListReq{}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		x.HandleErr(c, x.ErrorModel(constant.ParamError))
		return
	}

	token, err := service.OrderList(c, req)
	if err != nil {
		x.HandleErr(c, x.ErrorModel(constant.ServerError))
		return
	}

	x.Success(c, token)
}
