package model

import (
	"time"
)

// MessageUserGroupMapping 对应user_group_message_mappings表的数据模型
type MessageUserGroupMapping struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`
	GroupID    int       `json:"group_id" gorm:"column:group_id"`
	MessageID  int       `json:"message_id" gorm:"column:message_id"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (*MessageUserGroupMapping) TableName() string {
	return "message_user_group_mappings" // 表名指定为user_group_message_mappings
}
