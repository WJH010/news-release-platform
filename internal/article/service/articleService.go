package service

import (
	"context"
	"news-release/internal/article/model"
	"news-release/internal/article/repository"
)

// 服务接口，定义方法，接收 context.Context 和数据模型。
type ArticleService interface {
	ListArticle(ctx context.Context, page, pageSize int, articleTitle, articleType, releaseTime string, fieldID, isSelection, status int) ([]*model.Article, int64, error)
	GetArticleContent(ctx context.Context, articleID int) (*model.Article, error)
}

// 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type ArticleServiceImpl struct {
	articleRepo repository.ArticleRepository
}

// 创建服务实例
func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{articleRepo: articleRepo}
}

// 分页查询数据
func (s *ArticleServiceImpl) ListArticle(ctx context.Context, page, pageSize int, articleTitle, articleType, releaseTime string, fieldID, isSelection, status int) ([]*model.Article, int64, error) {
	return s.articleRepo.List(ctx, page, pageSize, articleTitle, articleType, releaseTime, fieldID, isSelection, status)
}

// 获取文章内容
func (s *ArticleServiceImpl) GetArticleContent(ctx context.Context, articleID int) (*model.Article, error) {
	return s.articleRepo.GetArticleContent(ctx, articleID)
}
