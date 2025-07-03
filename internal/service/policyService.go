package service

import (
	"context"
	"news-release/internal/model"
	"news-release/internal/repository"
)

// 服务接口，定义方法，接收 context.Context 和数据模型。
type PolicyService interface {
	ListPolicy(ctx context.Context, page, pageSize int, policyTitle string, fieldID int, is_selection int) ([]*model.Policy, int64, error)
	GetPolicyContent(ctx context.Context, policyID int) (*model.Policy, error)
}

// 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type PolicyServiceImpl struct {
	policyRepo repository.PolicyRepository
}

// 创建服务实例
func NewPolicyService(policyRepo repository.PolicyRepository) PolicyService {
	return &PolicyServiceImpl{policyRepo: policyRepo}
}

// 分页查询数据
func (s *PolicyServiceImpl) ListPolicy(ctx context.Context, page, pageSize int, policyTitle string, fieldID int, is_selection int) ([]*model.Policy, int64, error) {
	return s.policyRepo.List(ctx, page, pageSize, policyTitle, fieldID, is_selection)
}

// 获取政策内容
func (s *PolicyServiceImpl) GetPolicyContent(ctx context.Context, policyID int) (*model.Policy, error) {
	return s.policyRepo.GetPolicyContent(ctx, policyID)
}
