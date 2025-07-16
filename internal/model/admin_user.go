package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// AdminUser 管理系统用户模型
type AdminUser struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // 序列化时忽略密码（安全）
	Role      string    `json:"role" gorm:"default:'editor'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "user"
}

// Validate 验证数据
func (p *AdminUser) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
