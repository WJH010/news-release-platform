package repository

import (
	"context"
	"news-release/internal/model"

	"gorm.io/gorm"
)

// AdminUserRepository 管理用户仓库接口
type AdminUserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.AdminUser, error)
}

// AdminUserRepositoryImpl 实现
type AdminUserRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &AdminUserRepositoryImpl{db: db}
}

// GetByUsername 根据用户名查询管理用户
func (r *AdminUserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*model.AdminUser, error) {
	var admin model.AdminUser
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&admin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 账号不存在
		}
		return nil, result.Error // 数据库错误
	}
	return &admin, nil
}
