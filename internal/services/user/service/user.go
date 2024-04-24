package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xmicro/internal/constant"
	"xmicro/internal/nacos"
	"xmicro/internal/proto/pb"
	"xmicro/internal/services/user/dto"
	store_redis "xmicro/internal/store/redis"
	"xmicro/internal/utils/u_jwt"
	"xmicro/internal/x"
)

func UserLogin(c *gin.Context, req dto.LoginReq) (string, error) {

	token, err := u_jwt.CreateToken(
		u_jwt.Claims{
			Id:   "100",
			Name: req.Username,
		},
		[]byte(constant.JwtSalt))
	if err != nil {
		return token, x.ErrorModel(constant.ExternalServiceError)
	}

	if err = store_redis.SetRedisUserToken(c, "100", token); err != nil {
		return token, x.ErrorModel(constant.ExternalServiceError)
	}

	return token, nil
}

func UserOrders(c *gin.Context, req dto.UserOrdersReq) (string, error) {
	var (
		err     error
		rpcReq  = &pb.GetOrderByUserIdRequest{}
		rpcResp = &pb.GetOrderByUserIdResponse{}
	)add
	rpcReq.UserId = req.UserId
	if rpcResp, err = nacos.GetNextOrderClient().GetOrderByUserId(c, rpcReq); err != nil {
		return "", err
	}

	fmt.Println(rpcResp)

	return "", nil
}
