package repository

import (
	"context"
	"errors"
	"fmt"
	"news-release/internal/user/dto"
	"news-release/internal/user/model"
	"news-release/internal/utils"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	GetUserByOpenID(ctx context.Context, openid string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, userID int, updateFields map[string]interface{}) error
	UpdateSessionAndLoginTime(ctx context.Context, userID int, sessionKey string) error
	GetUserByID(ctx context.Context, userID int) (*dto.UserInfoResponse, error)
	ListAllUsers(ctx context.Context, page, pageSize int, req dto.ListUsersRequest) ([]*dto.ListUsersResponse, int64, error)
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
func (repo *UserRepositoryImpl) GetUserByOpenID(ctx context.Context, openid string) (*model.User, error) {
	var user model.User
	result := repo.db.WithContext(ctx).Where("openid = ?", openid).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", result.Error))
	}
	return &user, nil
}

// Create 创建新用户
func (repo *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	err := repo.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("创建用户失败: %w", err))
	}
	return err
}

// UpdateSessionAndLoginTime 登录时更新session_key和最后登录时间
func (repo *UserRepositoryImpl) UpdateSessionAndLoginTime(ctx context.Context, userID int, sessionKey string) error {
	result := repo.db.WithContext(ctx).Model(&model.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"session_key":     sessionKey,
			"last_login_time": time.Now(),
		})

	if result.Error != nil {
		return utils.NewSystemError(fmt.Errorf("更新登录信息失败: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return utils.NewSystemError(fmt.Errorf("更新登录信息异常，未更新任何数据: %w", result.Error))
	}

	return nil
}

// GetUserByID 获取用户信息
func (repo *UserRepositoryImpl) GetUserByID(ctx context.Context, userID int) (*dto.UserInfoResponse, error) {
	var user dto.UserInfoResponse
	query := repo.db.WithContext(ctx)

	result := query.Table("users u").
		Select(`u.nickname, u.avatar_url, u.name, u.gender AS gender_code,
				CASE
					WHEN gender = 'M' THEN
					'男'
					WHEN gender = 'F' THEN
					'女'
					ELSE
					'未知'
				END AS gender, u.phone_number, u.email, u.unit, u.department, u.position, u.industry, i.industry_name`).
		Joins("LEFT JOIN industries i ON u.industry = i.industry_code").
		Where("user_id = ?", userID).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, utils.NewSystemError(fmt.Errorf("查询用户信息失败: %v", result.Error))
	}
	return &user, nil
}

// Update 更新用户信息
func (repo *UserRepositoryImpl) Update(ctx context.Context, userID int, updateFields map[string]interface{}) error {

	result := repo.db.WithContext(ctx).Model(&model.User{}).
		Where("user_id = ?", userID).
		Updates(updateFields)

	if result.Error != nil {
		return utils.NewSystemError(fmt.Errorf("更新用户信息失败: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "更新用户信息失败，用户数据异常，请刷新页面后重试")
	}
	return nil
}

// ListAllUsers 分页查询用户列表
func (repo *UserRepositoryImpl) ListAllUsers(ctx context.Context, page, pageSize int, req dto.ListUsersRequest) ([]*dto.ListUsersResponse, int64, error) {
	var users []*dto.ListUsersResponse
	var total int64

	query := repo.db.WithContext(ctx).Table("users u").
		Select(`u.user_id, u.nickname, u.avatar_url, u.name, u.gender AS gender_code,
				CASE
					WHEN gender = 'M' THEN
					'男'
					WHEN gender = 'F' THEN
					'女'
					ELSE
					'未知'
				END AS gender,
				u.phone_number, u.email, u.unit, u.department, u.position, u.industry, i.industry_name, ur.role_name`).
		Joins("LEFT JOIN industries i ON u.industry = i.industry_code").
		Joins("LEFT JOIN user_role ur ON ur.role_code = u.role")

	// 拼接查询条件
	if req.Name != "" {
		query = query.Where("u.name LIKE ?", "%"+req.Name+"%")
	}
	if req.GenderCode != "" {
		query = query.Where("u.gender = ?", req.GenderCode)
	}
	if req.Unit != "" {
		query = query.Where("u.unit LIKE ?", "%"+req.Unit+"%")
	}
	if req.Department != "" {
		query = query.Where("u.department LIKE ?", "%"+req.Department+"%")
	}
	if req.Position != "" {
		query = query.Where("u.position LIKE ?", "%"+req.Position+"%")
	}
	if req.Industry != "" {
		query = query.Where("u.industry = ?", req.Industry)
	}

	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算用户总数失败: %w", err))
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("查询用户列表失败: %w", err))
	}

	return users, total, nil
}
