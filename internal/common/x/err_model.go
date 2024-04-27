package x

import (
	"xmicro/internal/common/constant"
)

type ServiceErrorModel struct {
	Code int32
	Msg  string
}

func (model *ServiceErrorModel) Error() string {
	return model.Msg
}

func ErrorModel(code int32) *ServiceErrorModel {
	return &ServiceErrorModel{Code: code, Msg: constant.MessageMap[code]}
}
