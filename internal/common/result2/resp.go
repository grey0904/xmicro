package result2

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"xmicro/internal/common/constant"
)

type Response struct {
	Code    int32       `json:"code"`
	Msg     string      `json:"msg"`
	TraceId string      `json:"tid,omitempty"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, d interface{}) {
	// TODO 日志记录
	resp(c, constant.Success, d, "success")
	return
}

func HandleErr(ctx *gin.Context, err error) {

	code := constant.SystemError
	msg := constant.MessageMap[constant.SystemError]
	resultData := err

	var serviceError *ServiceErrorModel
	if errors.As(err, &serviceError) {
		//logs.EF(ctx.Request.Context(), "手动抛出err: %+v", serviceError)
		code = serviceError.Code
		msg = constant.MessageMap[code]
		resultData = nil
	}

	//logs.EF(ctx.Request.Context(), "path = %s code = %d message = %s err = %+v", ctx.Request.URL, code, msg, err)
	resp(ctx, code, resultData, msg)
}

func resp(c *gin.Context, code int32, d interface{}, msg string) {

	span := trace.SpanFromContext(c.Request.Context())

	resp := Response{
		Code:    code,
		Msg:     msg,
		TraceId: span.SpanContext().TraceID().String(),
	}

	if d == nil {
		resp.Data = struct{}{}
	} else {
		resp.Data = d
	}

	c.JSON(http.StatusOK, resp)
}
