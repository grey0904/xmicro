package httperror

import "fmt"

type Http struct {
	Msg  string `json:"msg,omitempty"`
	Data string `json:"data,omitempty"`
	Code int    `json:"code"`
}

func (e Http) Error() string {
	return fmt.Sprintf("msg: %s,  data: %s", e.Msg, e.Data)
}

func NewHttpError(msg, data string, code int) Http {
	return Http{
		Msg:  msg,
		Data: data,
		Code: code,
	}
}
