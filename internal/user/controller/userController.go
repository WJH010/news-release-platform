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

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
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

	ctx.JSON(http.StatusOK, gin.H{
		"page":      req.Page,
		"page_size": req.PageSize,
		"total":     total,
		"data":      users,
	})
}

// BgLogin 后台登录
func (ctr *UserController) BgLogin(ctx *gin.Context) {
	// 绑定并验证请求参数
	var req dto.BgLoginRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	token, err := ctr.userService.BgLogin(ctx, req)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"token":   token,
	})
}

// CreateAdminUser 新增管理员
func (ctr *UserController) CreateAdminUser(ctx *gin.Context) {
	// 绑定并验证请求参数
	var req dto.CreateAdminRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取当前登录用户ID
	operator, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务新增管理员
	err = ctr.userService.CreateAdminUser(ctx, req, operator)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "新增管理员成功",
	})
}

// UpdateAdminUser 更新管理员
func (ctr *UserController) UpdateAdminUser(ctx *gin.Context) {
	// 从路径参数获取userID
	var urlReq dto.UserIDRequest
	if !utils.BindUrl(ctx, &urlReq) {
		return
	}

	// 绑定并验证请求参数
	var req dto.UpdateAdminRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取当前登录用户ID
	operator, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务更新管理员
	err = ctr.userService.UpdateAdminUser(ctx, urlReq.UserID, req, operator)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新管理员成功",
	})
}

// UpdateAdminStatus 更新管理员状态
func (ctr *UserController) UpdateAdminStatus(ctx *gin.Context) {
	// 从路径参数获取userID
	var urlReq dto.UserIDRequest
	if !utils.BindUrl(ctx, &urlReq) {
		return
	}

	// 绑定并验证请求参数
	var req dto.UpdateAdminStatusRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取当前登录用户ID
	operator, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务更新管理员状态
	err = ctr.userService.UpdateAdminStatus(ctx, urlReq.UserID, req.Operation, operator)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新管理员状态成功",
	})
}
