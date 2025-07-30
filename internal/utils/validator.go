package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators 向现有验证器注册自定义规则
func RegisterCustomValidators(v *validator.Validate) {
	// 时间格式验证
	err := v.RegisterValidation("time_format", validateTimeFormat)
	if err != nil {
		panic("注册时间格式验证失败: " + err.Error())
	}

	// 其他自定义规则...
}

// 时间格式验证实现
func validateTimeFormat(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	if fieldValue == "" {
		return true // 空值通过（配合omitempty使用）
	}
	// 定义允许的时间格式
	allowedFormats := []string{
		"2006-01-02",          // 日期格式
		"2006-01-02 15:04:05", // 日期时间格式
	}

	// 尝试用每种格式解析
	for _, format := range allowedFormats {
		if _, err := time.Parse(format, fieldValue); err == nil {
			return true
		}
	}

	return false
}
