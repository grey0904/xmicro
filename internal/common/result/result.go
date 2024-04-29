package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code int `json:"code"`
	Msg  any `json:"msg"`
}

// Fail err 最后自己封装一个
func Fail(ctx *gin.Context, err *Error) {
	ctx.JSON(http.StatusOK, Result{
		Code: err.Code,
		Msg:  err.Err.Error(),
	})
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Result{
		Code: OK,
		Msg:  data,
	})
}
