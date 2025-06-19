package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// News 结构体对应数据库中的new表
type News struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	NewTitle     string    `json:"new_title" gorm:"column:new_title;type:varchar(255);not null"`
	FieldName    string    `json:"field_name" gorm:"column:field_name"` // 关联字段
	ReleaseTime  time.Time `json:"release_time" gorm:"column:release_time;type:datetime"`
	BriefContent string    `json:"brief_content" gorm:"column:brief_content;type:text"`
	NewContent   string    `json:"new_content" gorm:"column:new_content;type:mediumtext"`
	Status       string    `json:"status" gorm:"column:status;type:varchar(255)"`
	IsSelection  int       `json:"is_selection" gorm:"column:is_selection;type:int"`
	FieldID      int       `json:"field_id" gorm:"column:field_id;type:int"`
	CreationTime time.Time `json:"creation_time" gorm:"column:creation_time;type:date"`
	UpdateTime   time.Time `json:"update_time" gorm:"column:update_time;type:date"`
}

// TableName 指定表名
func (News) TableName() string {
	return "new_items"
}

// Validate 验证数据
func (p *News) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
