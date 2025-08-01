package repository

import (
	"context"
	"errors"
	"fmt"
	"news-release/internal/user/model"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	GetUserByOpenID(ctx context.Context, openid string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdateSessionAndLoginTime(ctx context.Context, userID int, sessionKey string) error
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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

// UpdateSessionAndLoginTime 登录时更新session_key和最后登录时间
func (r *UserRepositoryImpl) UpdateSessionAndLoginTime(ctx context.Context, userID int, sessionKey string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"session_key":     sessionKey,
			"last_login_time": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("更新登录信息失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
