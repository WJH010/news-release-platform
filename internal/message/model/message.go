package model

import (
	"time"
)

// 消息相关常量定义
const (
	MessageTypeEvent = "EVENT" // 活动消息类型
)

// Message 对应messages表的数据模型
type Message struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`
	Title      string    `json:"title" gorm:"type:varchar(255);column:title"`
	Content    string    `json:"content" gorm:"type:mediumtext;column:content"`
	SendTime   time.Time `json:"send_time" gorm:"column:send_time"`
	Type       string    `json:"type" gorm:"type:varchar(50);column:type"`
	GroupID    int       `json:"group_id" gorm:"column:group_id"`               // 关联的用户组ID
	IsDeleted  string    `json:"is_deleted" gorm:"column:is_deleted;default:N"` // 软删除标志，默认值为N
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	// 关联字段
	TypeName string `json:"type_name" gorm:"column:type_name"` // 关联message_types表
}

// TableName 设置表名
func (*Message) TableName() string {
	return "messages" // 表名指定为messages
}
