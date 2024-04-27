package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
)

type UserController struct {
}

func (u *UserController) Register(ctx *gin.Context) {
	//接收参数
	var req pb.RegisterParams
	err2 := ctx.ShouldBindJSON(&req)
	if err2 != nil {
		common.Fail(ctx, biz.RequestDataError)
		return
	}
	response, err := rpc.UserClient.Register(context.TODO(), &req)
	if err != nil {
		common.Fail(ctx, msError.ToError(err))
		return
	}
	uid := response.Uid
	if len(uid) == 0 {
		common.Fail(ctx, biz.SqlError)
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
		common.Fail(ctx, biz.Fail)
		return
	}
	result := map[string]any{
		"token": token,
		"serverInfo": map[string]any{
			"host": config.Conf.Services["connector"].ClientHost,
			"port": config.Conf.Services["connector"].ClientPort,
		},
	}
	common.Success(ctx, result)
}
