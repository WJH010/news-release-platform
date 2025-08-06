package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"news-release/internal/user/model"
	"news-release/internal/utils"
)

// UserGroupRepository 用户群组仓库接口
type UserGroupRepository interface {
	// GetUserGroupByEventID 获取用户组信息
	GetUserGroupByEventID(ctx context.Context, eventID int) (*model.UserGroup, error)
	// AddUserToGroup 添加用户到用户组
	AddUserToGroup(ctx context.Context, userGroupMap *model.UserGroupMapping) error
}

// UserGroupRepositoryImpl 用户群组仓库实现
type UserGroupRepositoryImpl struct {
	db *gorm.DB
}

// NewUserGroupRepository 创建用户组仓库实例
func NewUserGroupRepository(db *gorm.DB) UserGroupRepository {
	return &UserGroupRepositoryImpl{db: db}
}

// GetUserGroupByEventID 根据用户组ID获取用户组信息
func (repo *UserGroupRepositoryImpl) GetUserGroupByEventID(ctx context.Context, eventID int) (*model.UserGroup, error) {
	var group model.UserGroup
	result := repo.db.WithContext(ctx).Where("event_id = ? AND id_deleted = ?", eventID, "N").First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", result.Error))
	}
	return &group, nil
}

// AddUserToGroup 将用户添加到用户组
func (repo *UserGroupRepositoryImpl) AddUserToGroup(ctx context.Context, userGroupMap *model.UserGroupMapping) error {
	err := repo.db.WithContext(ctx).Create(userGroupMap).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			// 如果是重复键错误，说明用户已在该用户组中
			return utils.NewBusinessError(utils.ErrCodeResourceExists, "用户已在该活动群组中")
		}
		return utils.NewSystemError(fmt.Errorf("添加用户到用户组失败: %w", err))
	}
	return err
}
