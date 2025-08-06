package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

// BindQuery Query参数绑定函数
func BindQuery(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindQuery(req); err != nil {
		var validationErr validator.ValidationErrors
		var numErr *strconv.NumError
		if errors.As(err, &validationErr) {
			// 转换为友好提示
			// msg := GetValidationErrorMsg(validationErr[0])
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			err = NewBusinessError(ErrCodeParamInvalid, msg.String())
			WrapErrorHandler(ctx, err)
		} else if errors.As(err, &numErr) {
			// 捕获数字转换错误（如非数字字符串传入int字段）
			err = NewBusinessError(ErrCodeParamInvalid, "参数类型错误，请检查参数格式")
			WrapErrorHandler(ctx, err)
		} else {
			err = NewBusinessError(ErrCodeParamBind, "参数绑定失败")
			WrapErrorHandler(ctx, err)
		}

		return false
	}
	return true
}

// BindJSON JSON格式请求体绑定函数
func BindJSON(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			// 转换为友好提示
			// msg := GetValidationErrorMsg(validationErr[0])
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			err = NewBusinessError(ErrCodeParamInvalid, msg.String())
			WrapErrorHandler(ctx, err)
		} else {
			err = NewBusinessError(ErrCodeParamBind, "参数绑定失败")
			WrapErrorHandler(ctx, err)
		}
		return false
	}
	return true
}

// BindUrl URL 路径参数绑定函数
func BindUrl(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindUri(req); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			// 转换为友好提示
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			err = NewBusinessError(ErrCodeParamInvalid, msg.String())
			WrapErrorHandler(ctx, err)
		} else {
			err = NewBusinessError(ErrCodeParamBind, "参数绑定失败")
			WrapErrorHandler(ctx, err)
		}
		return false
	}
	return true
}

// BindForm 表单参数绑定函数
func BindForm(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBind(req); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			// 转换为友好提示
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			err = NewBusinessError(ErrCodeParamInvalid, msg.String())
			WrapErrorHandler(ctx, err)
		} else {
			err = NewBusinessError(ErrCodeParamBind, "参数绑定失败")
			WrapErrorHandler(ctx, err)
		}
		return false
	}
	return true
}

// GetValidationErrorMsg 转换验证错误为友好提示
func GetValidationErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + "为必填参数"
	case "min":
		return err.Field() + "最小值为" + err.Param()
	case "max":
		return err.Field() + "最大值为" + err.Param()
	case "oneof":
		return err.Field() + "值不在有效范围内"
	case "numeric":
		return err.Field() + "必须为数字"
	case "time_format":
		return err.Field() + "日期格式错误"
	case "nickname":
		return err.Field() + "昵称格式错误，必须为2-20个字符，支持中英文、数字和下划线"
	case "real_name":
		return err.Field() + "姓名格式错误，必须为2-10个中文字符或4-20个英文字符"
	case "phone":
		return err.Field() + "手机号格式错误"
	case "email":
		return err.Field() + "邮箱格式错误"
	default:
		return err.Field() + "参数格式错误"
	}
}
