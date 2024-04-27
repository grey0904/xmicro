package dto

type LoginReq struct {
	Username string `json:"username" binding:"required"`
}

type UserOrdersReq struct {
	UserId int64 `json:"user_id" binding:"required"`
}
