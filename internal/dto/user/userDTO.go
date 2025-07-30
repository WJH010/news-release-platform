package user

type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}
