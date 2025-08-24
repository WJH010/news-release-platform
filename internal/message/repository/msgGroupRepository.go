package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"news-release/internal/message/dto"
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
	// DeleteUserMsgGroupMappings 删除用户-消息群组关联记录（软删除）
	DeleteUserMsgGroupMappings(ctx context.Context, msgGroupID int, userIDs []int, operateUser int) error
	// UpdateMsgGroup 更新消息群组信息
	UpdateMsgGroup(ctx context.Context, msgGroupID int, updateField map[string]interface{}) error
	// ListMsgGroups 分页查询消息群组
	ListMsgGroups(ctx context.Context, page int, pageSize int, groupName string, eventID int, queryScope string) ([]model.UserMessageGroup, int64, error)
	// ListGroupsUsers 查询指定群组的用户列表
	ListGroupsUsers(ctx context.Context, page int, pageSize int, msgGroupID int) ([]dto.ListGroupsUsersResponse, int64, error)
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

// CreateMsgGroup 创建消息群组
func (repo *MsgGroupRepositoryImpl) CreateMsgGroup(ctx context.Context, group *model.UserMessageGroup) error {
	err := repo.db.WithContext(ctx).Create(group).Error
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("创建消息群组失败: %v", err))
	}
	return nil
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

// DeleteUserMsgGroupMappings 删除用户-消息群组关联记录（软删除）
func (repo *MsgGroupRepositoryImpl) DeleteUserMsgGroupMappings(ctx context.Context, msgGroupID int, userIDs []int, operateUser int) error {
	if len(userIDs) == 0 {
		return nil
	}

	if err := repo.db.WithContext(ctx).Model(&model.UserMsgGroupMapping{}).
		Where("msg_group_id = ? AND user_id in (?) AND is_deleted = ?", msgGroupID, userIDs, "N").
		Updates(map[string]interface{}{
			"is_deleted":  "Y",
			"update_user": operateUser,
		}).Error; err != nil {
		return utils.NewSystemError(fmt.Errorf("批量删除用户-消息群组关联记录失败: %v", err))
	}

	return nil
}

// UpdateMsgGroup 更新消息群组信息
func (repo *MsgGroupRepositoryImpl) UpdateMsgGroup(ctx context.Context, msgGroupID int, updateField map[string]interface{}) error {
	err := repo.db.WithContext(ctx).Model(&model.UserMessageGroup{}).
		Where("id = ?", msgGroupID).
		Updates(updateField).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return utils.NewSystemError(fmt.Errorf("更新消息群组失败: %v", err))
	}
	return nil
}

// ListMsgGroups 分页查询消息群组
func (repo *MsgGroupRepositoryImpl) ListMsgGroups(ctx context.Context, page int, pageSize int, groupName string, eventID int, queryScope string) ([]model.UserMessageGroup, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	var groups []model.UserMessageGroup

	query := repo.db.WithContext(ctx).Table("user_message_groups umg").
		Select("umg.id, umg.group_name, umg.desc, umg.event_id, e.event_title, umg.include_all_user, umg.is_deleted").
		Joins("LEFT JOIN events e ON e.id = umg.event_id")
	// 拼接查询条件
	if groupName != "" {
		query = query.Where("umg.group_name LIKE ?", "%"+groupName+"%")
	}
	if eventID != 0 {
		query = query.Where("umg.event_id = ?", eventID)
	}
	if queryScope != "" {
		// 如果传入了查询范围，则添加查询条件
		// 如果传入了查询范围为DELETED，则查询已删除的群组
		if queryScope == utils.QueryScopeDeleted {
			query = query.Where("umg.is_deleted = ?", utils.DeletedFlagYes) // 查询已删除的文章
		}
		if queryScope == utils.QueryScopeAll {
			// 如果传入了查询范围为ALL，则查询所有群组，包括已删除和未删除的
		}
	} else {
		// 默认查询未删除群组
		query = query.Where("umg.is_deleted = ?", utils.DeletedFlagNo)
	}

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算总数时数据库查询失败: %v", err))
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&groups).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return groups, total, nil
}

// ListGroupsUsers 查询指定群组的用户列表
func (repo *MsgGroupRepositoryImpl) ListGroupsUsers(ctx context.Context, page int, pageSize int, msgGroupID int) ([]dto.ListGroupsUsersResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	var users []dto.ListGroupsUsersResponse
	query := repo.db.WithContext(ctx)

	query = query.Table("users u").
		Select(`
				u.user_id, u.nickname, u.name, u.gender AS gender_code, 
				CASE 
					WHEN u.gender = 'M' THEN '男' 
					WHEN gender = 'F' THEN '女' 
					ELSE '未知'
				END AS gender,
				u.phone_number, u.email, u.unit, u.department, u.position, 
				u.industry, i.industry_name, m.is_deleted"`).
		Joins("LEFT JOIN industries i ON u.industry = i.industry_code").
		Joins("JOIN user_msg_group_mappings m ON u.user_id = m.user_id").
		Where("m.msg_group_id = ? AND m.is_deleted = ?", msgGroupID, "N")

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算总数时数据库查询失败: %v", err))
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}
	return users, total, nil
}
