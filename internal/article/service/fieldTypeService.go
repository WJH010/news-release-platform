package service

import (
	"context"
	"news-release/internal/article/model"
	"news-release/internal/article/repository"
)

// FieldTypeService 服务接口，定义方法，接收 context.Context 和数据模型。
type FieldTypeService interface {
	GetFieldType(ctx context.Context) ([]*model.FieldType, error)
}

// FieldTypeServiceImpl 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type FieldTypeServiceImpl struct {
	fieldTypeRepo repository.FieldTypeRepository
}

// NewFieldTypeService 创建服务实例
func NewFieldTypeService(fieldTypeRepo repository.FieldTypeRepository) FieldTypeService {
	return &FieldTypeServiceImpl{fieldTypeRepo: fieldTypeRepo}
}

// GetFieldType 获取领域类型列表
func (svc *FieldTypeServiceImpl) GetFieldType(ctx context.Context) ([]*model.FieldType, error) {
	return svc.fieldTypeRepo.GetFieldType(ctx)
}
