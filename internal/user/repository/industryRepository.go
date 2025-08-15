package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"news-release/internal/user/model"
	"news-release/internal/utils"
)

type IndustryRepository interface {
	// ListIndustries 查询行业列表
	ListIndustries(ctx context.Context) ([]*model.Industries, error)
}

type IndustryRepositoryImpl struct {
	db *gorm.DB // 假设使用 GORM 或其他 ORM 库
}

// NewIndustryRepository 创建行业数据访问实例
func NewIndustryRepository(db *gorm.DB) IndustryRepository {
	return &IndustryRepositoryImpl{db: db}
}

// ListIndustries 查询行业列表
func (repo *IndustryRepositoryImpl) ListIndustries(ctx context.Context) ([]*model.Industries, error) {
	var industries []*model.Industries
	result := repo.db.WithContext(ctx).Find(&industries)
	err := result.Error

	if err != nil {
		return nil, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return industries, nil
}
