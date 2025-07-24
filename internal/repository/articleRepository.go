package repository

import (
	"context"
	"net/http"
	"news-release/internal/model"
	"news-release/internal/utils"

	"gorm.io/gorm"
)

// 数据访问接口，定义数据访问的方法集
type ArticleRepository interface {
	// 分页查询
	List(ctx context.Context, page, pageSize int, articleTitle, articleType, releaseTime string, fieldID, isSelection, status int) ([]*model.Article, int64, error)
	// 内容查询
	GetArticleContent(ctx context.Context, articleID int) (*model.Article, error)
}

// 实现接口的具体结构体
type ArticleRepositoryImpl struct {
	db *gorm.DB
}

// 创建数据访问实例
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &ArticleRepositoryImpl{db: db}
}

// 分页查询数据
func (r *ArticleRepositoryImpl) List(ctx context.Context, page, pageSize int, articleTitle, articleType, releaseTime string, fieldID, isSelection, status int) ([]*model.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var articles []*model.Article
	query := r.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("articles a").
		Select("a.article_id, a.article_title, a.article_type, a.release_time, a.brief_content, a.is_selection, f.field_name, a.cover_image_url").
		Joins("LEFT JOIN field_type f ON a.field_id = f.field_id")

	// 添加条件查询
	if releaseTime != "" {
		query = query.Where("a.release_time >= ?", releaseTime)
	}
	if articleTitle != "" {
		query = query.Where("a.article_title LIKE ?", "%"+articleTitle+"%")
	}
	if fieldID != 0 {
		query = query.Where("a.field_id = ?", fieldID)
	}
	if articleType != "" {
		query = query.Where("a.article_type = ?", articleType)
	}
	if isSelection != 0 {
		query = query.Where("a.is_selection = ?", isSelection)
	}
	if status != 0 {
		query = query.Where("a.status = ?", status)
	}

	// 按发布时间降序排列
	query = query.Order("a.release_time DESC")

	// 计算总数
	var total int64
	countQuery := *query // 复制查询对象，避免修改原始查询
	if err := countQuery.Count(&total).Error; err != nil {
		utils.HandleError(nil, err, http.StatusInternalServerError, 0, "计算总数时数据库查询失败")
		return nil, 0, err
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		utils.HandleError(nil, err, http.StatusInternalServerError, 0, "数据库查询失败")
		return nil, 0, err
	}

	return articles, total, nil
}

// 内容查询
func (r *ArticleRepositoryImpl) GetArticleContent(ctx context.Context, articleID int) (*model.Article, error) {
	var article model.Article

	result := r.db.WithContext(ctx).First(&article, articleID)
	err := result.Error

	// 查询文章内容
	if err != nil {
		return nil, err
	}

	return &article, nil
}
