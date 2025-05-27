package controller

import (
	"net/http"
	"news-release/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExampleController 控制器
type ExampleController struct {
	exampleService service.ExampleService
}

// NewExampleController 创建控制器实例
func NewExampleController(exampleService service.ExampleService) *ExampleController {
	return &ExampleController{exampleService: exampleService}
}

// ListExample 分页查询
func (c *ExampleController) ListExample(ctx *gin.Context) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	field1 := ctx.Query("field1")

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
	example, total, err := c.exampleService.ListExample(ctx, page, pageSize, field1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回分页结果
	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      example,
	})
}
