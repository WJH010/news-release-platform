package repository

import (
	"context"
	"net/http"
	"news-release/internal/model"
	"news-release/internal/utils"

	"gorm.io/gorm"
)

// NewsRepository 新闻仓库接口
type NewsRepository interface {
	//分页查询新闻列表
	GetNewsList(ctx context.Context, page, pageSize int, newsTitle string, fieldID int) ([]*model.News, int64, error)
	//获取新闻内容
	GetNewsContent(ctx context.Context, newsID int) (*model.News, error)
}

// NewsRepositoryImpl 新闻仓库实现
type NewsRepositoryImpl struct {
	// 数据库连接或其他依赖
	db *gorm.DB
}

// NewNewsRepository 创建新闻仓库实例
func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &NewsRepositoryImpl{db: db}
}

// 分页查询数据
func (r *NewsRepositoryImpl) GetNewsList(ctx context.Context, page, pageSize int, NewsTitle string, fieldID int) ([]*model.News, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var new []*model.News
	query := r.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("new_items p").
		Select("p.id, p.new_title, p.release_time, p.brief_content, f.field_name, p.list_image_url").
		Joins("LEFT JOIN field_type f ON p.field_id = f.field_id")

	// 添加条件查询
	if NewsTitle != "" {
		query = query.Where("p.new_title LIKE ?", "%"+NewsTitle+"%")
	}
	if fieldID != 0 {
		query = query.Where("p.field_id = ?", fieldID)
	}

	// 按发布时间降序排列
	query = query.Order("p.release_time DESC")

	// 计算总数
	var total int64
	countQuery := *query // 复制查询对象，避免修改原始查询
	if err := countQuery.Count(&total).Error; err != nil {
		utils.HandleError(nil, err, http.StatusInternalServerError, 0, "计算总数时数据库查询失败")
		return nil, 0, err
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&new).Error; err != nil {
		utils.HandleError(nil, err, http.StatusInternalServerError, 0, "数据库查询失败")
		return nil, 0, err
	}

	return new, total, nil
}

// 新闻内容查询
func (r *NewsRepositoryImpl) GetNewsContent(ctx context.Context, newsID int) (*model.News, error) {
	var new model.News

	result := r.db.WithContext(ctx).First(&new, newsID)
	err := result.Error

	// 查询新闻内容
	if err != nil {
		return nil, err
	}

	return &new, nil
}
