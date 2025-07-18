package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// User 数据模型
type User struct {
	UserID        int64     `json:"user_id" gorm:"primaryKey;column:user_id"`
	OpenID        string    `json:"openid" gorm:"column:openid"`
	UnionID       string    `json:"unionid" gorm:"column:unionid;default:NULL"`
	SessionKey    string    `json:"session_key" gorm:"column:session_key"`
	Nickname      string    `json:"nickname" gorm:"column:nickname"`
	AvatarURL     string    `json:"avatar_url" gorm:"column:avatar_url"`
	Gender        int       `json:"gender" gorm:"column:gender"`
	PhoneNumber   string    `json:"phone_number" gorm:"column:phone_number;default:NULL"`
	Email         string    `json:"email" gorm:"column:email"`
	Region        string    `json:"region" gorm:"column:region"`
	Status        int       `json:"status" gorm:"column:status;default:1"` // 默认=1，1：正常，2：禁用
	LastLoginTime time.Time `json:"last_login_time" gorm:"column:last_login_time;autoUpdateTime"`
	UserLevel     int       `json:"user_level" gorm:"column:user_level"`
	Password      string    `json:"password" gorm:"column:password"`
	Role          int       `json:"role" gorm:"column:role;default:1"` // 默认=1，1：普通用户，2：管理员
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime    time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

// Validate 验证数据
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
