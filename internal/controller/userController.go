package controller

import (
	"net/http"

	"news-release/internal/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// Login 微信登录接口
func (c *UserController) Login(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		utils.HandleError(ctx, nil, http.StatusBadRequest, 0, "缺少 code 参数")
		return
	}

	user, err := c.userService.Login(ctx, code)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "登录失败")
		return
	}

	if user == nil {
		utils.HandleError(ctx, nil, http.StatusInternalServerError, 0, "获取用户信息失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    user,
	})
}
