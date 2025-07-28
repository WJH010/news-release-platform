package user

import (
	"net/http"

	usersvc "news-release/internal/service/user"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService usersvc.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService usersvc.UserService) *UserController {
	return &UserController{userService: userService}
}

// Login 微信登录接口
func (c *UserController) Login(ctx *gin.Context) {
	// 获取请求参数
	var request struct {
		Code string `json:"code" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.HandleError(ctx, nil, http.StatusBadRequest, 0, "缺少 code 参数"+err.Error())
		return
	}

	token, err := c.userService.Login(ctx, request.Code)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "登录失败"+err.Error())
		return
	}

	if token == "" {
		utils.HandleError(ctx, nil, http.StatusInternalServerError, 0, "token生成异常")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"token":   token,
	})
}
