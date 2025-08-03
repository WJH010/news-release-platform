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
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "服务器内部错误，获取领域类型列表失败")
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": fieldType})
}
