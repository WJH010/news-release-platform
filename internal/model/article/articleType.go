package article

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// article数据模型
type ArticleType struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`
	TypeCode   string    `json:"type_code" gorm:"column:type_code"`
	TypeName   string    `json:"type_name" gorm:"column:type_name"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 设置表名
func (ArticleType) TableName() string {
	return "article_type"
}

// Validate 验证数据
func (a *ArticleType) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
