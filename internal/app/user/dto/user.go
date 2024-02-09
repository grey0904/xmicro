package dto

type LoginReq struct {
	Username string `json:"username" binding:"required"`
}
