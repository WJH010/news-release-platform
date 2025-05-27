package model

import (
	"github.com/go-playground/validator/v10"
)

// 数据模型
type Example struct {
	// 字段        类型       标签
	// json:控制结构体字段在 JSON 序列化 / 反序列化时的名称
	// gorm:控制数据库表字段的约束
	// validate:控制字段的验证规则
	ID     uint   `json:"id" gorm:"primaryKey"`
	FIELD1 string `json:"field1"`
}

// TableName 设置表名
// 当不指定表名时，GORM 默认使用结构体名作为表名
func (Example) TableName() string {
	return "example"
}

// Validate 验证用户数据
// 使用 go-playground/validator 包来进行数据验证
// 验证规则来源于结构体标签中的 validate
func (u *Example) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
