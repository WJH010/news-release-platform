package service

import (
	"context"
	"news-release/internal/user/model"
	"news-release/internal/user/repository"
)

// IndustryService 服务接口，定义行业相关的业务逻辑方法
type IndustryService interface {
	// ListIndustries 查询行业列表
	ListIndustries(ctx context.Context) ([]*model.Industries, error)
}

// IndustryServiceImpl 实现 IndustryService 接口，提供行业相关的业务逻辑
type IndustryServiceImpl struct {
	industryRepo repository.IndustryRepository // 行业数据访问接口
}

// NewIndustryService 创建服务实例
func NewIndustryService(industryRepo repository.IndustryRepository) IndustryService {
	return &IndustryServiceImpl{
		industryRepo: industryRepo,
	}
}

// ListIndustries 查询行业列表
func (svc *IndustryServiceImpl) ListIndustries(ctx context.Context) ([]*model.Industries, error) {
	// 调用数据访问层的查询方法
	return svc.industryRepo.ListIndustries(ctx)
}
