package service

import (
	"context"
	"news-release/internal/model"
	"news-release/internal/repository"
)

// ExampleService 服务接口，定义方法，接收 context.Context 和数据模型。
type ExampleService interface {
	ListExample(ctx context.Context, page, pageSize int, field1 string) ([]*model.Example, int64, error)
}

// ExampleServiceImpl 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type ExampleServiceImpl struct {
	exampleRepo repository.ExampleRepository
}

// NewExampleService 创建服务实例
func NewExampleService(exampleRepo repository.ExampleRepository) ExampleService {
	return &ExampleServiceImpl{exampleRepo: exampleRepo}
}

// ListExample 分页查询数据
func (s *ExampleServiceImpl) ListExample(ctx context.Context, page, pageSize int, field1 string) ([]*model.Example, int64, error) {
	return s.exampleRepo.List(ctx, page, pageSize, field1)
}
