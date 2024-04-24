package dto

type OrderListReq struct {
	UserId string `json:"user_id" binding:"required"`
}
