package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Code    int    `json:"code"`              // 业务错误码
	Message string `json:"message"`           // 用户友好的错误提示
	Details string `json:"details,omitempty"` // 详细错误信息（可选）
}

/*
*
用于处理错误。
使用示例：

	if err != nil {
	    HandleError(ctx, err, http.StatusInternalServerError, 10000, "错误描述")
	    return
	}

*
*/
func HandleError(ctx *gin.Context, err error, statusCode int, code int, message string) {
	logrus.WithError(err).
		WithField("status", statusCode).
		WithField("errorCode", code).
		Error(message)

	if ctx != nil {
		ctx.JSON(statusCode, ErrorResponse{
			Code:    code,
			Message: message,
			Details: err.Error(),
		})
		ctx.Abort() // 终止后续处理
	}
}

// Query参数绑定函数
func BindQuery(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindQuery(req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			// 转换为友好提示
			// msg := GetValidationErrorMsg(validationErr[0])
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamInvalid, msg.String())
		} else if _, ok := err.(*strconv.NumError); ok {
			// 捕获数字转换错误（如非数字字符串传入int字段）
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamInvalid, "查询参数格式错误")
		} else {
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamBind, "参数绑定失败")
		}

		return false
	}
	return true
}

// JSON格式请求体绑定函数
func BindJSON(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			// 转换为友好提示
			// msg := GetValidationErrorMsg(validationErr[0])
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamInvalid, msg.String())
		} else {
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamBind, "参数绑定失败")
		}
		return false
	}
	return true
}

// URL 路径参数绑定函数
func BindUrl(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindUri(req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			// 转换为友好提示
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamInvalid, msg.String())
		} else {
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamBind, "参数绑定失败")
		}
		return false
	}
	return true
}

// 表单参数绑定函数
func BindForm(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBind(req); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			// 转换为友好提示
			var msg strings.Builder
			for _, e := range validationErr {
				msg.WriteString(GetValidationErrorMsg(e) + "; ")
			}
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamInvalid, msg.String())
		} else {
			HandleError(ctx, err, http.StatusBadRequest, ErrCodeParamBind, "参数绑定失败")
		}
		return false
	}
	return true
}

// 转换验证错误为友好提示
func GetValidationErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + "为必填参数"
	case "min":
		return err.Field() + "最小值为" + err.Param()
	case "max":
		return err.Field() + "最大值为" + err.Param()
	case "oneof":
		return err.Field() + "必须为" + err.Param() + "中的一个"
	case "numeric":
		return err.Field() + "必须为数字"
	case "time_format":
		return err.Field() + "日期格式错误"
	default:
		return err.Field() + "参数格式错误"
	}
}
