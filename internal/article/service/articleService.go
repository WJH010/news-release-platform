package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"news-release/internal/article/model"
	"news-release/internal/article/repository"
	db "news-release/internal/database"
	filerepo "news-release/internal/file/repository"
	"news-release/internal/utils"
)

// ArticleService 服务接口，定义方法，接收 context.Context 和数据模型。
type ArticleService interface {
	// ListArticle 分页查询文章列表
	ListArticle(ctx context.Context, page, pageSize int, articleTitle string, articleType string, releaseTime string, fieldType string, isSelection int) ([]*model.Article, int64, error)
	// GetArticleContent 获取文章内容
	GetArticleContent(ctx context.Context, articleID int) (*model.Article, error)
	// CreateArticle 创建文章
	CreateArticle(ctx context.Context, article *model.Article, imageIDList []int) error
}

// ArticleServiceImpl 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type ArticleServiceImpl struct {
	articleRepo repository.ArticleRepository
	fileRepo    filerepo.FileRepository
}

// NewArticleService 创建服务实例
func NewArticleService(articleRepo repository.ArticleRepository, fileRepo filerepo.FileRepository) ArticleService {
	return &ArticleServiceImpl{articleRepo: articleRepo, fileRepo: fileRepo}
}

// ListArticle 分页查询数据
func (svc *ArticleServiceImpl) ListArticle(ctx context.Context, page, pageSize int, articleTitle string, articleType string, releaseTime string, fieldType string, isSelection int) ([]*model.Article, int64, error) {
	return svc.articleRepo.List(ctx, page, pageSize, articleTitle, articleType, releaseTime, fieldType, isSelection)
}

// GetArticleContent 获取文章内容
func (svc *ArticleServiceImpl) GetArticleContent(ctx context.Context, articleID int) (*model.Article, error) {
	return svc.articleRepo.GetArticleContent(ctx, articleID)
}

// CreateArticle 创建文章
func (svc *ArticleServiceImpl) CreateArticle(ctx context.Context, article *model.Article, imageIDList []int) error {
	// 检查是否存在重复标题的文章
	existingArticle, err := svc.articleRepo.GetArticleByTitle(ctx, article.ArticleTitle)
	if err != nil {
		return err
	}
	if existingArticle != nil {
		return utils.NewBusinessError(utils.ErrCodeResourceExists, "已存在同名文章，请修改标题后重试")
	}

	// 开启事务
	tx := db.GetDB().Begin()
	if tx.Error != nil {
		return utils.NewSystemError(fmt.Errorf("开启事务失败: %w", tx.Error))
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logrus.Panic("事务回滚，发生异常: ", r)
		}
	}()

	// 创建文章
	if err := svc.articleRepo.CreateArticle(ctx, tx, article); err != nil {
		tx.Rollback()
		return err
	}

	// 如果有图片，更新images表的biz_id和biz_type
	if len(imageIDList) > 0 {
		if err := svc.fileRepo.BatchUpdateImageBizID(ctx, tx, imageIDList, article.ArticleID, utils.TypeArticle); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.NewSystemError(fmt.Errorf("提交事务失败: %w", err))
	}
	return nil
}
