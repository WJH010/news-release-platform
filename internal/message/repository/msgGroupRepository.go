package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"news-release/internal/message/model"
	"news-release/internal/utils"
)

// MsgGroupRepository 消息群组数据访问接口
type MsgGroupRepository interface {
	// ExecTransaction 执行事务
	ExecTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	// CreateMsgGroup 创建消息群组
	CreateMsgGroup(ctx context.Context, group *model.UserMessageGroup) error
	// GetMsgGroupByID 根据ID获取消息群组
	GetMsgGroupByID(ctx context.Context, msgGroupID int) (*model.UserMessageGroup, error)
	// GetExistingMappings 查询指定群组中已存在的用户关联记录
	GetExistingMappings(ctx context.Context, groupID int, userIDs []int) (map[int]model.UserMsgGroupMapping, error)
	// CreateUserMsgGroupMappings 批量创建用户-消息群组关联记录
	CreateUserMsgGroupMappings(ctx context.Context, tx *gorm.DB, mappings []model.UserMsgGroupMapping) error
	// RecoverUserMsgGroupMappings 批量恢复用户-消息群组关联记录
	RecoverUserMsgGroupMappings(ctx context.Context, tx *gorm.DB, msgGroupID int, userIDs []int, lastReadMsgID int, operateUser int) error
}

// MsgGroupRepositoryImpl 实现消息群组数据访问接口的具体结构体
type MsgGroupRepositoryImpl struct {
	db          *gorm.DB
	messageRepo MessageRepository
}

// NewMsgGroupRepository 创建消息群组数据访问实例
func NewMsgGroupRepository(db *gorm.DB, messageRepo MessageRepository) MsgGroupRepository {
	return &MsgGroupRepositoryImpl{db: db, messageRepo: messageRepo}
}

// ExecTransaction 实现事务执行（使用 GORM 的 Transaction 方法）
func (repo *MsgGroupRepositoryImpl) ExecTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return repo.db.WithContext(ctx).Transaction(fn)
}

// 工具函数：从需要恢复的记录中提取用户ID
func extractUserIDs(mappings []model.UserMsgGroupMapping) []int {
	var ids []int
	for _, m := range mappings {
		ids = append(ids, m.UserID)
	}
	return ids
}

// GetMsgGroupByID 根据ID获取消息群组
func (repo *MsgGroupRepositoryImpl) GetMsgGroupByID(ctx context.Context, msgGroupID int) (*model.UserMessageGroup, error) {
	var group model.UserMessageGroup
	err := repo.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", msgGroupID, "N").First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, utils.NewSystemError(fmt.Errorf("查询消息群组失败: %v", err))
	}
	return &group, nil
}

// GetExistingMappings 查询指定群组中已存在的用户关联记录
func (repo *MsgGroupRepositoryImpl) GetExistingMappings(ctx context.Context, groupID int, userIDs []int) (map[int]model.UserMsgGroupMapping, error) {
	var mappings []model.UserMsgGroupMapping
	if err := repo.db.WithContext(ctx).
		Where("msg_group_id = ? AND user_id IN (?)", groupID, userIDs).
		Find(&mappings).Error; err != nil {
		return nil, err
	}

	// 转换为 map[userID]mapping，便于快速查询
	result := make(map[int]model.UserMsgGroupMapping, len(mappings))
	for _, m := range mappings {
		result[m.UserID] = m
	}
	return result, nil
}

// CreateUserMsgGroupMappings 批量创建用户-消息群组关联记录
func (repo *MsgGroupRepositoryImpl) CreateUserMsgGroupMappings(ctx context.Context, tx *gorm.DB, mappings []model.UserMsgGroupMapping) error {
	if len(mappings) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).Create(&mappings).Error; err != nil {
		return utils.NewSystemError(fmt.Errorf("批量创建用户-消息群组关联记录失败: %v", err))
	}
	return nil
}

// RecoverUserMsgGroupMappings 批量恢复用户-消息群组关联记录
func (repo *MsgGroupRepositoryImpl) RecoverUserMsgGroupMappings(ctx context.Context, tx *gorm.DB, msgGroupID int, userIDs []int, lastReadMsgID int, operateUser int) error {
	if len(userIDs) == 0 {
		return nil
	}

	if err := tx.WithContext(ctx).Model(&model.UserMsgGroupMapping{}).
		Where("msg_group_id = ? AND user_id in (?)", msgGroupID, userIDs).
		Updates(map[string]interface{}{
			"is_deleted":       "N",
			"last_read_msg_id": lastReadMsgID,
			"join_msg_id":      lastReadMsgID,
			"update_user":      operateUser,
		}).Error; err != nil {
		return utils.NewSystemError(fmt.Errorf("批量恢复用户-消息群组关联记录失败: %v", err))
	}

	return nil
}

// CreateMsgGroup 创建消息群组
func (repo *MsgGroupRepositoryImpl) CreateMsgGroup(ctx context.Context, group *model.UserMessageGroup) error {
	err := repo.db.WithContext(ctx).Create(group).Error
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("创建消息群组失败: %v", err))
	}
	return nil
}
