package controller

import (
	"net/http"
	"news-release/internal/article/dto"
	"news-release/internal/article/model"
	"news-release/internal/article/service"
	"news-release/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
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
	article, total, err := ctr.articleService.ListArticle(ctx, page, pageSize, req.ArticleTitle, req.ArticleType, req.ReleaseTime, req.FieldType, req.IsSelection)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
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
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
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

// CreateArticle 创建文章
func (ctr *ArticleController) CreateArticle(ctx *gin.Context) {
	// 初始化参数结构体并绑定请求体
	var req dto.CreateArticleRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 构造文章对象
	article := &model.Article{
		ArticleTitle:   req.ArticleTitle,
		ArticleType:    req.ArticleType,
		BriefContent:   req.BriefContent,
		ArticleContent: req.ArticleContent,
		IsSelection:    req.IsSelection,
		FieldType:      req.FieldType,
		CoverImageURL:  req.CoverImageURL,
		ArticleSource:  req.ArticleSource,
		ReleaseTime:    time.Now(),
		CreateUser:     userID,
		UpdateUser:     userID,
	}

	// 调用服务层
	if err := ctr.articleService.CreateArticle(ctx, article, req.ImageIDList); err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文章创建成功",
		"data": gin.H{
			"article_id": article.ArticleID,
		},
	})
}

// UpdateArticle 处理更新文章的请求
func (ctr *ArticleController) UpdateArticle(ctx *gin.Context) {
	// 从URL获取文章ID
	var urlReq dto.ArticleContentRequest
	if !utils.BindUrl(ctx, &urlReq) {
		return
	}

	// 从请求体获取更新参数
	var req dto.UpdateArticleRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取当前用户ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层更新文章
	err = ctr.articleService.UpdateArticle(ctx, urlReq.ArticleID, req, userID)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文章更新成功",
	})
}

// DeleteArticle 处理删除文章的请求
func (ctr *ArticleController) DeleteArticle(ctx *gin.Context) {
	// 从URL获取文章ID
	var req dto.ArticleContentRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 获取当前用户ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层删除文章
	err = ctr.articleService.DeleteArticle(ctx, req.ArticleID, userID)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文章删除成功",
	})
}
