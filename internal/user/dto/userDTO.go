package dto

type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	Nickname    *string `json:"nickname" binding:"omitempty,nickname"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,url"`
	Name        *string `json:"name" binding:"omitempty,real_name"`
	Gender      *int    `json:"gender" binding:"omitempty,oneof=1 2 3"` // 1: 男, 2: 女, 3: 未知
	PhoneNumber *string `json:"phone_number" binding:"omitempty,phone"`
	Email       *string `json:"email" binding:"omitempty,email"`
	Unit        *string `json:"unit" binding:"omitempty"`
	Department  *string `json:"department" binding:"omitempty"`
	Position    *string `json:"position" binding:"omitempty"`
	Industry    *string `json:"industry" binding:"omitempty"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Unit        string `json:"unit"`
	Department  string `json:"department"`
	Position    string `json:"position"`
	Industry    string `json:"industry"`
}
