package service

import (
	"github.com/gin-gonic/gin"
	"xmicro/internal/app/order/dto"
)

func OrderList(c *gin.Context, req dto.OrderListReq) (string, error) {

	//app.Db.Where("user_id = ?", 1).Find(&userOrders)

	return "", nil
}
