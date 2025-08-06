package controller

import (
	"net/http"
	"news-release/internal/article/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

// FieldTypeController 控制器
type FieldTypeController struct {
	fieldTypeService service.FieldTypeService
}

// NewFieldTypeController 创建控制器实例
func NewFieldTypeController(fieldTyprService service.FieldTypeService) *FieldTypeController {
	return &FieldTypeController{fieldTypeService: fieldTyprService}
}

// GetFieldType 获取政策内容
func (ctr *FieldTypeController) GetFieldType(ctx *gin.Context) {

	// 调用服务层
	fieldType, err := ctr.fieldTypeService.GetFieldType(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": fieldType})
}
