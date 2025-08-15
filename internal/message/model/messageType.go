package model

import (
	"time"
)

// MessageType 对应 news_platform 平台下 message_types 数据表的数据模型
type MessageType struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`             // 主键ID
	TypeCode   string    `json:"type_code" gorm:"not null;column:type_code"` // 消息类型代码
	TypeName   string    `json:"type_name" gorm:"column:type_name"`          // 消息类型描述
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`      // 数据创建时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`      // 数据最后更新时间
}

// TableName 设置当前模型对应的数据库表名
func (*MessageType) TableName() string {
	return "message_types"
}
