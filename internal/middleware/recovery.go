package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				logrus.WithField("panic", err).Error("发生panic")

				// 返回500错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "服务器内部错误",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}
