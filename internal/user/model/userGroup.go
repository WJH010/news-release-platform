package model

import (
	"time"
)

// UserGroup 对应user_groups表的数据模型
type UserGroup struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	GroupCode   string    `json:"group_code" gorm:"type:varchar(20);column:group_code"`
	GroupName   string    `json:"group_name" gorm:"type:varchar(50);column:group_name"`
	Description string    `json:"description" gorm:"type:varchar(255);column:description"`
	IsDeleted   string    `json:"is_deleted" gorm:"column:is_deleted;default:N"` // 软删除标志，默认值为N
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (*UserGroup) TableName() string {
	return "user_groups" // 表名指定为user_groups
}
