package service

import (
	"context"
	"errors"
	"fmt"
	"news-release/internal/model"
	"news-release/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// AdminUserService 管理用户服务接口
type AdminUserService interface {
	Login(ctx context.Context, username, password string) (*model.AdminUser, error)
}

// AdminUserServiceImpl 实现
type AdminUserServiceImpl struct {
	adminRepo repository.AdminUserRepository
}

func NewAdminUserService(adminRepo repository.AdminUserRepository) AdminUserService {
	return &AdminUserServiceImpl{adminRepo: adminRepo}
}

// Login 管理系统登录逻辑（验证账号密码）
func (s *AdminUserServiceImpl) Login(ctx context.Context, username, password string) (*model.AdminUser, error) {
	// 1. 查询用户是否存在
	admin, err := s.adminRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err // 数据库错误
	}
	if admin == nil {
		return nil, errors.New("账号不存在")
	}

	// 2. 验证密码（bcrypt比对哈希）
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errors.New("密码错误")
		}
		return nil, fmt.Errorf("密码验证出错: %w", err)
	}

	return admin, nil
}
