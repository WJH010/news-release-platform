package dto

type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}
