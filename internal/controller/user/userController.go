package user

import (
	"net/http"

	userdto "news-release/internal/dto/user"
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
	// 初始化参数结构体并绑定查询参数
	var req userdto.WxLoginRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	token, err := c.userService.Login(ctx, req.Code)
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
