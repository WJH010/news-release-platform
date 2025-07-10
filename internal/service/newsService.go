package service

import (
	"context"
	"news-release/internal/model"
	"news-release/internal/repository"
)

// NewsService 新闻服务接口
type NewsService interface {
	GetNewsList(ctx context.Context, page, pageSize int, newsTitle string, fieldID int, is_selection int) ([]*model.News, int64, error)
	GetNewsContent(ctx context.Context, newsID int) (*model.News, error)
}

// NewsServiceImpl 新闻服务实现
type NewsServiceImpl struct {
	newsRepo repository.NewsRepository
	// 可以注入其他依赖
}

// NewNewsService 创建新闻服务实例
func NewNewsService(newsRepo repository.NewsRepository) NewsService {
	return &NewsServiceImpl{
		newsRepo: newsRepo,
	}
}

func (s *NewsServiceImpl) GetNewsList(ctx context.Context, page, pageSize int, newsTitle string, fieldID int, is_selection int) ([]*model.News, int64, error) {
	return s.newsRepo.GetNewsList(ctx, page, pageSize, newsTitle, fieldID, is_selection)
}

// 获取新闻内容
func (s *NewsServiceImpl) GetNewsContent(ctx context.Context, newsID int) (*model.News, error) {
	return s.newsRepo.GetNewsContent(ctx, newsID)
}
