package controller

import (
	"fmt"
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
func (ctr *UserController) Login(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.WxLoginRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	token, err := ctr.userService.Login(ctx, req.Code)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	if token == "" {
		err = utils.NewSystemError(fmt.Errorf("token生成异常"))
		utils.WrapErrorHandler(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"token":   token,
	})
}

// UpdateUserInfo 更新用户信息接口
func (ctr *UserController) UpdateUserInfo(ctx *gin.Context) {
	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 绑定并验证请求参数
	var req dto.UserUpdateRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 调用服务更新用户信息
	err = ctr.userService.UpdateUserInfo(ctx, userID, req)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// GetUserInfo 获取用户信息接口
func (ctr *UserController) GetUserInfo(ctx *gin.Context) {
	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务获取用户信息
	user, err := ctr.userService.GetUserByID(ctx, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	result := dto.UserInfoResponse{
		Nickname:    user.Nickname,
		AvatarURL:   user.AvatarURL,
		Name:        user.Name,
		GenderCode:  user.Gender,
		Gender:      map[string]string{"M": "男", "F": "女", "U": "未知"}[user.Gender],
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Unit:        user.Unit,
		Department:  user.Department,
		Position:    user.Position,
		Industry:    user.Industry,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
