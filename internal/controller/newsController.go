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

// NewsListResponse 新闻列表响应结构体
type NewsListResponse struct {
	ID           int       `json:"id"`
	NewTitle     string    `json:"new_title"`
	FieldName    string    `json:"field_name"`
	ReleaseTime  time.Time `json:"release_time"`
	BriefContent string    `json:"brief_content"`
	ListImageURL string    `json:"list_image_url"`
}

// NewsContentResponse 新闻内容响应结构体
type NewsContentResponse struct {
	ID          int       `json:"id"`
	NewTitle    string    `json:"new_title"`
	ReleaseTime time.Time `json:"release_time"`
	NewContent  string    `json:"new_content"`
}

// NewsController 新闻控制器
type NewsController struct {
	newsService service.NewsService
}

// NewNewsController 创建新闻控制器实例
func NewNewsController(newsService service.NewsService) *NewsController {
	return &NewsController{
		newsService: newsService,
	}
}

func (p *NewsController) GetNewsList(ctx *gin.Context) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	newTitle := ctx.Query("newTitle")
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
	news, total, err := p.newsService.GetNewsList(ctx, page, pageSize, newTitle, fieldID)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []NewsListResponse
	for _, p := range news {
		result = append(result, NewsListResponse{
			ID:           p.ID,
			NewTitle:     p.NewTitle,
			FieldName:    p.FieldName,
			ReleaseTime:  p.ReleaseTime,
			BriefContent: p.BriefContent,
			ListImageURL: p.ListImageURL,
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
func (p *NewsController) GetNewsContent(ctx *gin.Context) {
	// 获取主键
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusBadRequest, 0, "无效的新闻ID")
		return
	}

	// 调用服务层
	news, err := p.newsService.GetNewsContent(ctx, int(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(ctx, err, http.StatusNotFound, 0, "新闻不存在(id="+idStr+")")
			return
		}
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "获取新闻内容失败")
		return
	}

	result := NewsContentResponse{
		ID:            news.ID,
		NewTitle:   news.NewTitle,
		ReleaseTime:   news.ReleaseTime,
		NewContent: news.NewContent,
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
