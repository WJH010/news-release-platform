package repository

import (
	"context"
	"news-release/internal/model"

	"gorm.io/gorm"
)

// ExampleRepository 数据访问接口，定义数据访问的方法集（Create，Read，Update，Delete）
type ExampleRepository interface {
	// 分页查询
	List(ctx context.Context, page, pageSize int, field1 string) ([]*model.Example, int64, error)
}

// ExampleRepositoryImpl 实现接口的具体结构体
type ExampleRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建数据访问实例
func NewExampleRepository(db *gorm.DB) ExampleRepository {
	return &ExampleRepositoryImpl{db: db}
}

// List 分页查询数据
func (r *ExampleRepositoryImpl) List(ctx context.Context, page, pageSize int, field1 string) ([]*model.Example, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var users []*model.Example
	query := r.db.WithContext(ctx)

	// 添加条件查询
	if field1 != "" {
		query = query.Where("field1 LIKE ?", "%"+field1+"%")
	}

	// 计算总数
	var total int64
	if err := query.Model(&model.Example{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
