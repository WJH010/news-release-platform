package controller

import (
	"net/http"

	"news-release/internal/user/dto"
	"news-release/internal/user/service"
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
	// 初始化参数结构体并绑定查询参数
	var req dto.WxLoginRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	token, err := c.userService.Login(ctx, req.Code)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "登录失败")
		return
	}

	if token == "" {
		utils.HandleError(ctx, nil, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "token生成异常")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"token":   token,
	})
}

// UpdateUserInfo 更新用户信息接口
func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	// 获取userID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeAuthFailed, "获取用户ID失败")
		return
	}

	// 绑定并验证请求参数
	var req dto.UserUpdateRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 调用服务更新用户信息
	err = c.userService.UpdateUserInfo(ctx, userID, req)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "更新用户信息失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}
