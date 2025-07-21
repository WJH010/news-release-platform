package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Notice struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Title       string     `json:"title" gorm:"type:varchar(255)"`
	Content     string     `json:"content" gorm:"type:text;not null"`
	ReleaseTime *time.Time `json:"release_time" gorm:"column:release_time"`
	Status      int        `json:"status" gorm:"type:int;not null;default:1"` // 默认=1，1：有效
	CreateTime  *time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  *time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (Notice) TableName() string {
	return "notices"
}

// Validate 验证数据
func (n *Notice) Validate() error {
	validate := validator.New()
	return validate.Struct(n)
}
