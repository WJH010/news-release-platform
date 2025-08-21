package controller

import (
	"github.com/gin-gonic/gin"
	"news-release/internal/user/dto"
	"news-release/internal/user/model"
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

// CreateIndustry 创建行业
func (ctr *IndustryController) CreateIndustry(ctx *gin.Context) {
	// 绑定请求体到 DTO
	var req dto.CreateIndustryRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	industry := &model.Industries{
		IndustryCode: req.IndustryCode,
		IndustryName: req.IndustryName,
	}

	// 调用服务层创建行业
	err := ctr.industryService.CreateIndustry(ctx, industry)
	if err != nil {
		// 处理异常
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回成功响应
	ctx.JSON(201, gin.H{"message": "行业创建成功"})
}

// UpdateIndustry 更新行业信息
func (ctr *IndustryController) UpdateIndustry(ctx *gin.Context) {
	// 绑定url参数获取行业ID
	var urlReq dto.IndustryUrlID
	if !utils.BindUrl(ctx, &urlReq) {
		return
	}

	// 绑定请求体到 DTO
	var req dto.UpdateIndustryRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 调用服务层更新行业信息
	err := ctr.industryService.UpdateIndustry(ctx, urlReq.ID, req)
	if err != nil {
		// 处理异常
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回成功响应
	ctx.JSON(200, gin.H{"message": "行业更新成功"})
}

// DeleteIndustry 删除行业
func (ctr *IndustryController) DeleteIndustry(ctx *gin.Context) {
	// 绑定url参数获取行业ID
	var urlReq dto.IndustryUrlID
	if !utils.BindUrl(ctx, &urlReq) {
		return
	}

	// 调用服务层删除行业
	err := ctr.industryService.DeleteIndustry(ctx, urlReq.ID)
	if err != nil {
		// 处理异常
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回成功响应
	ctx.JSON(200, gin.H{"message": "行业删除成功"})
}
