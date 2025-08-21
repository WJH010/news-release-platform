package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Code    int    `json:"code"`              // 业务错误码
	Message string `json:"message"`           // 用户友好的错误提示
	Details string `json:"details,omitempty"` // 详细错误信息（可选）
}

// HandleError
// 处理错误并返回统一格式的JSON响应
//- ctx: Gin上下文
//- err: 错误对象
//- statusCode: HTTP状态码
//- code: 业务错误码
//- message: 用户友好的错误消息

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

// WrapErrorHandler 封装错误处理逻辑，简化Controller层代码
// 注意：调用后需手动添加 return 终止当前函数，避免后续代码执行
func WrapErrorHandler(ctx *gin.Context, err error) {
	// 处理业务错误
	if bizErr, ok := GetBusinessError(err); ok {
		HandleError(ctx, err, http.StatusBadRequest, bizErr.Code, bizErr.Msg)
		return
	}

	// 处理系统错误
	var sysErr *SystemError
	if errors.As(err, &sysErr) {
		HandleError(ctx, sysErr.Err, http.StatusInternalServerError, ErrCodeServerInternalError, "服务器内部错误")
		return
	}

	// 处理未知错误
	HandleError(ctx, err, http.StatusInternalServerError, ErrCodeServerInternalError, "未知服务器错误")
}

// IsUniqueConstraintError 判断是否为唯一索引冲突错误
// 返回值：(是否为唯一冲突, 冲突字段名或索引名)
func IsUniqueConstraintError(err error) (bool, string) {
	// 先判断是否是GORM的错误（如记录未找到等，但唯一冲突通常是底层错误）
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, ""
	}

	// 提取底层错误（可能被多次包装）
	var underlyingErr error
	if errors.As(err, &underlyingErr) {
		err = underlyingErr
	}

	// 适配MySQL
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1062 { // MySQL唯一索引冲突错误码
			return true, parseMySQLUniqueField(mysqlErr.Message)
		}
	}

	return false, ""
}

// 解析MySQL错误信息中的冲突字段
func parseMySQLUniqueField(msg string) string {
	// 错误信息格式示例："Duplicate entry 'test' for key 'users.username'"
	// 提取 "users.username" 中的 "username" 部分
	start := strings.LastIndex(msg, ".")
	if start == -1 {
		return "unknown"
	}
	end := strings.LastIndex(msg, "'")
	if end == -1 || end <= start {
		return "unknown"
	}
	return msg[start+1 : end]
}
