package article

import (
	"net/http"
	articlersvc "news-release/internal/service/article"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

// 控制器
type FieldTypeController struct {
	fieldTyprService articlersvc.FieldTypeService
}

// 创建控制器实例
func NewFieldTypeController(fieldTyprService articlersvc.FieldTypeService) *FieldTypeController {
	return &FieldTypeController{fieldTyprService: fieldTyprService}
}

// 获取政策内容
func (f *FieldTypeController) GetFieldType(ctx *gin.Context) {

	// 调用服务层
	fieldType, err := f.fieldTyprService.GetFieldType(ctx)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "获取领域类型列表失败")
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": fieldType})
}
