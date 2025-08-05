package repository

import (
	"context"
	"fmt"
	"news-release/internal/article/model"

	"gorm.io/gorm"
)

// ArticleRepository 数据访问接口，定义数据访问的方法集
type ArticleRepository interface {
	// List 分页查询
	List(ctx context.Context, page, pageSize int, articleTitle string, articleType string, releaseTime string, fieldType string, isSelection int) ([]*model.Article, int64, error)
	// GetArticleContent 内容查询
	GetArticleContent(ctx context.Context, articleID int) (*model.Article, error)
}

// ArticleRepositoryImpl 实现接口的具体结构体
type ArticleRepositoryImpl struct {
	db *gorm.DB
}

// NewArticleRepository 创建数据访问实例
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &ArticleRepositoryImpl{db: db}
}

// List 分页查询数据
func (repo *ArticleRepositoryImpl) List(ctx context.Context, page, pageSize int, articleTitle string, articleType string, releaseTime string, fieldType string, isSelection int) ([]*model.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var articles []*model.Article
	query := repo.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("articles a").
		Select("a.article_id, a.article_title, a.article_type, a.release_time, a.brief_content, a.is_selection, f.field_name, a.cover_image_url,a.article_source,at.type_name").
		Joins("LEFT JOIN field_types f ON a.field_type = f.field_code").
		Joins("LEFT JOIN article_types at ON a.article_type = at.type_code").
		Where("a.is_deleted = ?", "N")

	// 添加条件查询
	if releaseTime != "" {
		query = query.Where("a.release_time >= ?", releaseTime)
	}
	if articleTitle != "" {
		query = query.Where("a.article_title LIKE ?", "%"+articleTitle+"%")
	}
	if fieldType != "" {
		query = query.Where("a.field_type = ?", fieldType)
	}
	if articleType != "" {
		query = query.Where("a.article_type = ?", articleType)
	}
	if isSelection != 0 {
		query = query.Where("a.is_selection = ?", isSelection)
	}

	// 按发布时间降序排列
	query = query.Order("a.release_time DESC")

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算总数时数据库查询失败: %v", err)
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("数据库查询失败: %v", err)
	}

	return articles, total, nil
}

// GetArticleContent 内容查询
func (repo *ArticleRepositoryImpl) GetArticleContent(ctx context.Context, articleID int) (*model.Article, error) {
	var article model.Article

	result := repo.db.WithContext(ctx).First(&article, articleID)
	err := result.Error

	// 查询文章内容
	if err != nil {
		return nil, err
	}

	return &article, nil
}
