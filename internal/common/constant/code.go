package constant

const (
	Success     = int32(1000)
	SystemError = int32(1001)

	NeedUsername         = int32(1002)
	ParamError           = int32(1003)
	InvalidCredentials   = int32(1004)
	AccessDenied         = int32(1005)
	InvalidInput         = int32(1006)
	ResourceNotFound     = int32(1007)
	NetworkError         = int32(1008)
	ExternalServiceError = int32(1009)
	ServerError          = int32(1010)
	ResourceConflict     = int32(1011)
	TimeoutError         = int32(1012)
	OperationFailed      = int32(1013)
)

var MessageMap = map[int32]string{
	Success:              "success",
	SystemError:          "系统内部错误",
	NeedUsername:         "用户名不能为空",
	ParamError:           "参数错误",
	InvalidCredentials:   "无效的凭证",
	AccessDenied:         "访问被拒绝",
	InvalidInput:         "无效的输入",
	ResourceNotFound:     "资源未找到",
	NetworkError:         "网络错误",
	ExternalServiceError: "外部服务错误",
	ServerError:          "服务器错误",
	ResourceConflict:     "资源冲突",
	TimeoutError:         "操作超时",
	OperationFailed:      "操作失败",
}
