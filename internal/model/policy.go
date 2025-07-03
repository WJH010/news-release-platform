package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// policy数据模型
type Policy struct {
	ID            int       `json:"id" gorm:"primaryKey;column:id"`
	PolicyTitle   string    `json:"policy_title" gorm:"not null;column:policy_title"`
	FieldName     string    `json:"field_name" gorm:"column:field_name"` // 关联字段
	ReleaseTime   time.Time `json:"release_time" gorm:"column:release_time"`
	BriefContent  string    `json:"brief_content" gorm:"type:text;column:brief_content"`
	PolicyContent string    `json:"policy_content" gorm:"type:mediumtext;column:policy_content"`
	Status        int       `json:"status" gorm:"column:status;default:1"` // 默认=1，1：有效
	IsSelection   int       `json:"is_selection" gorm:"default:2;column:is_selection"` // 默认=2，1：精选，2：非精选
	FieldID       int       `json:"field_id" gorm:"column:field_id"`
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"update_time" gorm:"column:update_time"`
}

// TableName 设置表名
func (Policy) TableName() string {
	return "policy_items"
}

// Validate 验证数据
func (p *Policy) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
