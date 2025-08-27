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
		Nickname:     user.Nickname,
		AvatarURL:    user.AvatarURL,
		Name:         user.Name,
		GenderCode:   user.Gender,
		Gender:       map[string]string{"M": "男", "F": "女", "U": "未知"}[user.Gender],
		PhoneNumber:  user.PhoneNumber,
		Email:        user.Email,
		Unit:         user.Unit,
		Department:   user.Department,
		Position:     user.Position,
		Industry:     user.Industry,
		IndustryName: user.IndustryName,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// ListAllUsers 列出所有用户接口（管理员权限）
func (ctr *UserController) ListAllUsers(ctx *gin.Context) {
	// 绑定并验证请求参数
	var req dto.ListUsersRequest
	if !utils.BindQuery(ctx, &req) {
		return
	}

	// 设置默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 调用服务获取用户列表
	users, total, err := ctr.userService.ListAllUsers(ctx, req.Page, req.PageSize, req)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 构建响应数据
	var userResponses []dto.ListUsersResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.ListUsersResponse{
			UserID:       user.UserID,
			Nickname:     user.Nickname,
			AvatarURL:    user.AvatarURL,
			Name:         user.Name,
			GenderCode:   user.Gender,
			Gender:       map[string]string{"M": "男", "F": "女", "U": "未知"}[user.Gender],
			PhoneNumber:  user.PhoneNumber,
			Email:        user.Email,
			Unit:         user.Unit,
			Department:   user.Department,
			Position:     user.Position,
			Industry:     user.Industry,
			IndustryName: user.IndustryName,
			RoleName:     user.RoleName,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"page":      req.Page,
		"page_size": req.PageSize,
		"total":     total,
		"data":      userResponses,
	})
}

// TestLogin 测试用接口，直接返回token
func (ctr *UserController) TestLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"token":   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ5MjkxOTYsImlhdCI6MTc1NjI4OTE5Niwib3BlbmlkIjoib3AtOUl2cDZROWhhTEpqRDdIWU15TDJWMTNqOCIsInVzZXJfcm9sZSI6MiwidXNlcmlkIjo4fQ._g3jb63kMYQOU_RPaD-TBISb_dtioZJ9qekdZ3BbMiA",
	})
}
