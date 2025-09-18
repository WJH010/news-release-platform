package controller

import (
	"news-release/internal/user/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserRoleController struct {
	userRoleService service.UserRoleService // 用户角色服务接口
}

// NewUserRoleController 创建用户角色控制器实例
func NewUserRoleController(userRoleService service.UserRoleService) *UserRoleController {
	return &UserRoleController{
		userRoleService: userRoleService,
	}
}

// List 获取所有用户角色
func (ctr *UserRoleController) List(ctx *gin.Context) {
	list, err := ctr.userRoleService.List(ctx)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	// 返回成功响应
	ctx.JSON(200, gin.H{"data": list})
}
