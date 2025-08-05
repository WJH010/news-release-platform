package model

import (
	"time"
)

// MessageUserMapping 对应user_message_mappings表的数据模型
type MessageUserMapping struct {
	ID         int64     `json:"id" gorm:"primaryKey;column:id"`
	UserID     int       `json:"user_id" gorm:"column:user_id"`
	MessageID  int       `json:"message_id" gorm:"column:message_id"`
	IsRead     string    `json:"is_read" gorm:"type:varchar(2);column:is_read;default:N"`
	ReadTime   time.Time `json:"read_time" gorm:"column:read_time"`
	IsDeleted  string    `json:"is_deleted" gorm:"column:is_deleted;default:N"` // 软删除标志，默认值为N
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (*MessageUserMapping) TableName() string {
	return "message_user_mappings" // 表名指定为user_message_mappings
}
