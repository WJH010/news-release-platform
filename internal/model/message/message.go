package message

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Message 对应messages表的数据模型
type Message struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`
	Title      string    `json:"title" gorm:"type:varchar(255);column:title"`
	Content    string    `json:"content" gorm:"type:mediumtext;column:content"`
	SendTime   time.Time `json:"send_time" gorm:"column:send_time"`
	Type       string    `json:"type" gorm:"type:varchar(50);column:type"`
	Status     int       `json:"status" gorm:"column:status"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (Message) TableName() string {
	return "messages" // 表名指定为messages
}

// Validate 验证数据
func (m *Message) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}
