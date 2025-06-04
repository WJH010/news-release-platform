package repository

import (
	"context"
	"net/http"
	"news-release/internal/model"
	"news-release/internal/utils"

	"gorm.io/gorm"
)

// 数据访问接口，定义数据访问的方法集
type PolicyRepository interface {
	// 分页查询
	List(ctx context.Context, page, pageSize int, policyTitle string, fieldID int) ([]*model.Policy, int64, error)
}

// 实现接口的具体结构体
type PolicyRepositoryImpl struct {
	db *gorm.DB
}

// 创建数据访问实例
func NewPolicyRepository(db *gorm.DB) PolicyRepository {
	return &PolicyRepositoryImpl{db: db}
}

// 分页查询数据
func (r *PolicyRepositoryImpl) List(ctx context.Context, page, pageSize int, policyTitle string, fieldID int) ([]*model.Policy, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var policy []*model.Policy
	query := r.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("policy_items p").
		Select("p.id, p.policy_title, p.release_time, p.brief_content, f.field_name").
		Joins("LEFT JOIN field_type f ON p.field_id = f.field_id")

	// 添加条件查询
	if policyTitle != "" {
		query = query.Where("p.policy_title LIKE ?", "%"+policyTitle+"%")
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
		utils.HandleError(nil, err, http.StatusInternalServerError, 0, "数据库查询失败")
		return nil, 0, err
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&policy).Error; err != nil {
		utils.HandleError(nil, err, http.StatusInternalServerError, 0, "数据库查询失败")
		return nil, 0, err
	}

	return policy, total, nil
}
