package service

import (
	"github.com/gin-gonic/gin"
	dto "xmicro/internal/app/user/dto"
	"xmicro/internal/constant"
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
