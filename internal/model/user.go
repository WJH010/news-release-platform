package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// User 用户模型
type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	OpenID    string    `json:"openid" gorm:"unique;not null"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// Validate 验证数据
func (p *User) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
