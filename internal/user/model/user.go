package model

import (
	"time"
)

// User 数据模型
type User struct {
	UserID        int       `json:"user_id" gorm:"primaryKey;column:user_id"`
	OpenID        string    `json:"openid" gorm:"column:openid"`
	UnionID       string    `json:"unionid" gorm:"column:unionid;default:NULL"`
	SessionKey    string    `json:"session_key" gorm:"column:session_key"`
	Nickname      string    `json:"nickname" gorm:"column:nickname"`
	AvatarURL     string    `json:"avatar_url" gorm:"column:avatar_url"`
	Name          string    `json:"name" gorm:"column:name"`
	Gender        string    `json:"gender" gorm:"column:gender"` // M: 男, F: 女, U: 未知
	PhoneNumber   string    `json:"phone_number" gorm:"column:phone_number;default:NULL"`
	Email         string    `json:"email" gorm:"column:email"`
	Region        string    `json:"region" gorm:"column:region"`
	Status        int       `json:"status" gorm:"column:status;default:1"` // 默认=1，1：正常，2：禁用
	LastLoginTime time.Time `json:"last_login_time" gorm:"column:last_login_time;autoUpdateTime"`
	UserLevel     int       `json:"user_level" gorm:"column:user_level"`
	Password      string    `json:"password" gorm:"column:password"`
	Role          int       `json:"role" gorm:"column:role;default:1"` // 默认=1，1：普通用户，2：管理员
	Unit          string    `json:"unit" gorm:"column:unit;default:NULL"`
	Department    string    `json:"department" gorm:"column:department;default:NULL"`
	Position      string    `json:"position" gorm:"column:position;default:NULL"`
	Industry      int       `json:"industry" gorm:"column:industry;default:NULL"`
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime    time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	// 关联字段
	IndustryName string `json:"industry_name" gorm:"column:industry_name"` // 行业名称，关联 industries 表
	RoleName     string `json:"role_name" gorm:"column:role_name"`
}

// TableName 设置表名
func (*User) TableName() string {
	return "users"
}
