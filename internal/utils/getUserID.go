package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserID(ctx *gin.Context) (int, error) {
	// 获取userID
	userID, exists := ctx.Get("userid")
	if !exists {
		return 0, fmt.Errorf("获取用户ID失败")
	}
	// 类型转换
	uid, ok := userID.(int)
	if !ok {
		return 0, fmt.Errorf("用户ID类型错误")
	}
	return uid, nil
}
