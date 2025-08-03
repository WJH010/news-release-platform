package repository

import (
	"context"
	"news-release/internal/article/model"

	"gorm.io/gorm"
)

// FieldTypeRepository 数据访问接口，定义数据访问的方法集
type FieldTypeRepository interface {
	// GetFieldType 获取领域类型列表
	GetFieldType(ctx context.Context) ([]*model.FieldType, error)
}

// FieldTypeRepositoryImpl 实现接口的具体结构体
type FieldTypeRepositoryImpl struct {
	db *gorm.DB
}

// NewFieldTypeRepository 创建数据访问实例
func NewFieldTypeRepository(db *gorm.DB) FieldTypeRepository {
	return &FieldTypeRepositoryImpl{db: db}
}

// GetFieldType 获取领域类型列表
func (repo *FieldTypeRepositoryImpl) GetFieldType(ctx context.Context) ([]*model.FieldType, error) {
	var fieldType []*model.FieldType

	result := repo.db.WithContext(ctx).Find(&fieldType)
	err := result.Error

	if err != nil {
		return nil, err
	}

	return fieldType, nil
}
