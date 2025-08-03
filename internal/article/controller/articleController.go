package controller

import (
	"errors"
	"fmt"
	"net/http"
	"news-release/internal/article/dto"
	"news-release/internal/article/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArticleController 控制器
type ArticleController struct {
	articleService service.ArticleService
}

// NewArticleController 创建控制器实例
func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

// ListArticle 分页查询
func (ctr *ArticleController) ListArticle(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.ArticleListRequest
	if !utils.BindQuery(ctx, &req) {
		return
	}

	// page 默认1
	page := req.Page
	if page == 0 {
		page = 1
	}

	// pageSize 默认10
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	// 调用服务层
	article, total, err := ctr.articleService.ListArticle(ctx, page, pageSize, req.ArticleTitle, req.ArticleType, req.ReleaseTime, req.FieldID, req.IsSelection, req.Status)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "服务器内部错误，获取文章列表失败")
		return
	}

	var result []dto.ArticleListResponse
	for _, a := range article {
		result = append(result, dto.ArticleListResponse{
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

// GetArticleContent 获取文章内容
func (ctr *ArticleController) GetArticleContent(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.ArticleContentRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 调用服务层
	article, err := ctr.articleService.GetArticleContent(ctx, req.ArticleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("文章不存在(id=%d)", req.ArticleID)
			utils.HandleError(ctx, err, http.StatusNotFound, utils.ErrCodeResourceNotFound, msg)
			return
		}
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "服务器内部错误，获取文章内容失败")
		return
	}

	result := dto.ArticleContentResponse{
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
