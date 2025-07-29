package article

import (
	"errors"
	"net/http"
	articlersvc "news-release/internal/service/article"
	"news-release/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArticleListResponse 文章列表响应结构体
type ArticleListResponse struct {
	ArticleID       int       `json:"article_id"`
	ArticleTitle    string    `json:"article_title"`
	ArticleTypeCode string    `json:"article_type_code"`
	ArticleType     string    `json:"article_type"`
	FieldName       string    `json:"field_name"`
	ReleaseTime     time.Time `json:"release_time"`
	BriefContent    string    `json:"brief_content"`
	IsSelection     int       `json:"is_selection"`
	CoverImageURL   string    `json:"cover_image_url"`
	ArticleSource   string    `json:"article_source"`
}

// ArticleContentResponse 文章内容响应结构体
type ArticleContentResponse struct {
	ArticleID       int       `json:"article_id"`
	ArticleTitle    string    `json:"article_title"`
	FieldName       string    `json:"field_name"`
	ReleaseTime     time.Time `json:"release_time"`
	ArticleContent  string    `json:"article_content"`
	ArticleTypeCode string    `json:"article_type_code"`
	ArticleType     string    `json:"article_type"`
	ArticleSource   string    `json:"article_source"`
}

// 控制器
type ArticleController struct {
	articleService articlersvc.ArticleService
}

// 创建控制器实例
func NewArticleController(articleService articlersvc.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

// 分页查询
func (a *ArticleController) ListArticle(ctx *gin.Context) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	articleTitle := ctx.Query("article_title")
	fieldIDStr := ctx.Query("field_id")
	isSelectionStr := ctx.Query("is_selection")
	articleType := ctx.Query("article_type")
	releaseTime := ctx.Query("release_time")
	statusStr := ctx.Query("status")
	var fieldID int
	var isSelection int
	var status int

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

	// 转换 status 参数
	if statusStr != "" {
		var err error
		status, err = strconv.Atoi(statusStr)

		if err != nil {
			utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "status格式转换错误")
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
	article, total, err := a.articleService.ListArticle(ctx, page, pageSize, articleTitle, articleType, releaseTime, fieldID, isSelection, status)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []ArticleListResponse
	for _, a := range article {
		result = append(result, ArticleListResponse{
			ArticleID:       a.ArticleID,
			ArticleTitle:    a.ArticleTitle,
			ArticleTypeCode: a.ArticleType,
			ArticleType:     a.TypeName,
			FieldName:       a.FieldName,
			ReleaseTime:     a.ReleaseTime,
			BriefContent:    a.BriefContent,
			IsSelection:     a.IsSelection,
			CoverImageURL:   a.CoverImageURL,
			ArticleSource:   a.ArticleSource,
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

// 获取文章内容
func (p *ArticleController) GetArticleContent(ctx *gin.Context) {
	// 获取主键
	idStr := ctx.Param("articleID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusBadRequest, 0, "无效的文章ID")
		return
	}

	// 调用服务层
	article, err := p.articleService.GetArticleContent(ctx, int(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(ctx, err, http.StatusNotFound, 0, "文章不存在(id="+idStr+")")
			return
		}
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "获取文章内容失败")
		return
	}

	result := ArticleContentResponse{
		ArticleID:       article.ArticleID,
		ArticleTitle:    article.ArticleTitle,
		FieldName:       article.FieldName,
		ReleaseTime:     article.ReleaseTime,
		ArticleContent:  article.ArticleContent,
		ArticleTypeCode: article.ArticleType,
		ArticleType:     article.TypeName,
		ArticleSource:   article.ArticleSource,
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
