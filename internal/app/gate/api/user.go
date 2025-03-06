package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"xmicro/internal/app/user/pb"
	"xmicro/internal/common/config/center"
	"xmicro/internal/common/jwts"
	"xmicro/internal/common/logs"
	"xmicro/internal/common/result"
	"xmicro/internal/common/rpc/discovery"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) Register(ctx *gin.Context) {
	//接收参数
	var req pb.RegisterRequest
	err2 := ctx.ShouldBindJSON(&req)
	if err2 != nil {
		result.Fail(ctx, result.RequestDataError)
		return
	}
	response, err := discovery.UserClient.Register(context.TODO(), &req)
	if err != nil {
		result.Fail(ctx, result.ToError(err))
		return
	}
	uid := response.Uid
	if len(uid) == 0 {
		result.Fail(ctx, result.SqlError)
		return
	}
	logs.Info("uid:%s", uid)
	//gen token by uid jwt  A.B.C A部分头（定义加密算法） B部分 存储数据  C部分签名
	claims := jwts.CustomClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token, err := jwts.GenToken(&claims, config.Conf.Jwt.Secret)
	if err != nil {
		logs.Error("Register jwt gen token err:%v", err)
		result.Fail(ctx, result.RequestFail)
		return
	}
	resp := map[string]any{
		"token": token,
		"serverInfo": map[string]any{
			"host": config.Conf.Services["connector"].ClientHost,
			"port": config.Conf.Services["connector"].ClientPort,
		},
	}
	result.Success(ctx, resp)
}
