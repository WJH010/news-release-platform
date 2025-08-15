package controller

import (
	"github.com/gin-gonic/gin"
	"news-release/internal/user/dto"
	"news-release/internal/user/service"
	"news-release/internal/utils"
)

type IndustryController struct {
	industryService service.IndustryService // 行业服务接口
}

// NewIndustryController 创建行业控制器实例
func NewIndustryController(industryService service.IndustryService) *IndustryController {
	return &IndustryController{
		industryService: industryService,
	}
}

// ListIndustries 查询行业列表
func (ctr *IndustryController) ListIndustries(ctx *gin.Context) {
	// 调用服务层查询行业列表
	industries, err := ctr.industryService.ListIndustries(ctx)
	if err != nil {
		// 处理异常
		utils.WrapErrorHandler(ctx, err)
		return
	}
	var list []dto.ListIndustriesResponse
	for _, industry := range industries {
		// 将行业数据转换为响应格式
		list = append(list, dto.ListIndustriesResponse{
			ID:           industry.ID,
			IndustryCode: industry.IndustryCode,
			IndustryName: industry.IndustryName,
		})
	}

	// 返回成功响应
	ctx.JSON(200, gin.H{"data": list})
}
