package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// UserGroup 对应user_groups表的数据模型
type UserGroup struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	GroupCode   string    `json:"group_code" gorm:"type:varchar(20);column:group_code"`
	GroupName   string    `json:"group_name" gorm:"type:varchar(50);column:group_name"`
	Description string    `json:"description" gorm:"type:varchar(255);column:description"`
	Status      int       `json:"status" gorm:"column:status;default:1"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (UserGroup) TableName() string {
	return "user_groups" // 表名指定为user_groups
}

// Validate 验证数据
func (u *UserGroup) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
