package model

import (
	"time"
)

// UserGroupMapping 对应user_group_mappings表的数据模型
type UserGroupMapping struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`
	UserID     int       `json:"user_id" gorm:"column:user_id"`
	GroupID    int       `json:"group_id" gorm:"column:group_id"`
	IsDeleted  string    `json:"is_deleted" gorm:"column:is_deleted;default:N"` // 软删除标志，默认值为N
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (*UserGroupMapping) TableName() string {
	return "user_group_mappings" // 表名指定为user_group_mappings
}
