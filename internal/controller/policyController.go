package controller

import (
	"errors"
	"net/http"
	"news-release/internal/service"
	"news-release/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PolicyListResponse 政策列表响应结构体
type PolicyListResponse struct {
	ID           int       `json:"id"`
	PolicyTitle  string    `json:"policy_title"`
	FieldName    string    `json:"field_name"`
	ReleaseTime  time.Time `json:"release_time"`
	BriefContent string    `json:"brief_content"`
	IsSelection  int       `json:"is_selection"`
}

// PolicyContentResponse 政策内容响应结构体
type PolicyContentResponse struct {
	ID            int       `json:"id"`
	PolicyTitle   string    `json:"policy_title"`
	ReleaseTime   time.Time `json:"release_time"`
	PolicyContent string    `json:"policy_content"`
}

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
	isSelectionStr := ctx.Query("is_selection")
	var fieldID int
	var isSelection int

	// 转换 fieldID 参数
	if fieldIDStr != "" {
		var err error
		fieldID, err = strconv.Atoi(fieldIDStr)

		if err != nil {
			utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "fieldID格式转换错误")
			return
		}
	}

	// 转换 isSelection 参数
	if isSelectionStr != "" {
		var err error
		isSelection, err = strconv.Atoi(isSelectionStr)

		if err != nil {
			utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "isSelection格式转换错误")
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
	policy, total, err := p.policyService.ListPolicy(ctx, page, pageSize, policyTitle, fieldID, isSelection)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []PolicyListResponse
	for _, p := range policy {
		result = append(result, PolicyListResponse{
			ID:           p.ID,
			PolicyTitle:  p.PolicyTitle,
			FieldName:    p.FieldName,
			ReleaseTime:  p.ReleaseTime,
			BriefContent: p.BriefContent,
			IsSelection:  p.IsSelection,
		})
	}

	// 返回分页结果
	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      result,
	})
}

// 获取政策内容
func (p *PolicyController) GetPolicyContent(ctx *gin.Context) {
	// 获取主键
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusBadRequest, 0, "无效的政策ID")
		return
	}

	// 调用服务层
	policy, err := p.policyService.GetPolicyContent(ctx, int(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(ctx, err, http.StatusNotFound, 0, "政策不存在(id="+idStr+")")
			return
		}
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "获取政策内容失败")
		return
	}

	result := PolicyContentResponse{
		ID:            policy.ID,
		PolicyTitle:   policy.PolicyTitle,
		ReleaseTime:   policy.ReleaseTime,
		PolicyContent: policy.PolicyContent,
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
