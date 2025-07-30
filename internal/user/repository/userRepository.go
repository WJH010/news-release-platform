package repository

import (
	"context"
	"news-release/internal/user/model"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	GetUserByOpenID(ctx context.Context, openid string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

// UserRepositoryImpl 用户仓库实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// GetUserByOpenID 根据 openid 获取用户信息
func (r *UserRepositoryImpl) GetUserByOpenID(ctx context.Context, openid string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("openid = ?", openid).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create 创建新用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
