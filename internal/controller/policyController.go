package controller

import (
	"net/http"
	"news-release/internal/service"
	"news-release/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 控制器
type PolicyController struct {
	policyService service.PolicyService
}

// 创建控制器实例
func NewPolicyController(policyService service.PolicyService) *PolicyController {
	return &PolicyController{policyService: policyService}
}

// 分页查询
func (p *PolicyController) ListPolicy(ctx *gin.Context) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	policyTitle := ctx.Query("policyTitle")
	fieldIDStr := ctx.Query("fieldID")
	var fieldID int

	// 转换 fieldID 参数
	if fieldIDStr != "" {
		var err error
		fieldID, err = strconv.Atoi(fieldIDStr)

		if err != nil {
			utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "fieldID格式转换错误")
			return
		}
	}

	// 转换参数类型
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用服务层
	policy, total, err := p.policyService.ListPolicy(ctx, page, pageSize, policyTitle, fieldID)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	// 返回分页结果
	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      policy,
	})
}
