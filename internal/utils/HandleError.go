package utils

import (
	"github.com/gin-gonic/gin"
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
