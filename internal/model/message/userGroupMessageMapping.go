package message

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// UserGroupMessageMapping 对应user_group_message_mappings表的数据模型
type UserGroupMessageMapping struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`
	GroupID    int       `json:"group_id" gorm:"column:group_id"`
	MessageID  int       `json:"message_id" gorm:"column:message_id"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (UserGroupMessageMapping) TableName() string {
	return "user_group_message_mappings" // 表名指定为user_group_message_mappings
}

// Validate 验证数据
func (u *UserGroupMessageMapping) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
