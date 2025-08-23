package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"news-release/internal/message/dto"
	"news-release/internal/message/model"
	"news-release/internal/utils"
)

// MessageRepository 数据访问接口，定义数据访问的方法集
type MessageRepository interface {
	// GetMessageContent 内容查询
	GetMessageContent(ctx context.Context, messageID int) (*model.Message, error)
	// GetUnreadMessageCount 获取未读消息数
	GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error)
	// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
	// MarkAllMessagesAsRead(ctx context.Context, userID int) error
	// ListMessageGroupsByUserID 查询用户消息群组列表
	ListMessageGroupsByUserID(ctx context.Context, page, pageSize int, userID int, typeCode string) ([]*dto.MessageGroupDTO, int64, error)
	// ListMsgByGroups 分页查询分组内消息列表
	ListMsgByGroups(ctx context.Context, page, pageSize int, groupID int) ([]*dto.ListMessageDTO, int64, error)
	// MarkAsReadByGroup 按分组更新消息为已读
	MarkAsReadByGroup(ctx context.Context, userID int, groupID int)
	// 权限校验查询，确保普通用户只能查看自己的消息
	CheckUserMsgPermission(ctx context.Context, userID int, groupID int) error
}

// MessageRepositoryImpl 实现接口的具体结构体
type MessageRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageRepository 创建数据访问实例
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

// GetMessageContent 内容查询
func (repo *MessageRepositoryImpl) GetMessageContent(ctx context.Context, messageID int) (*model.Message, error) {
	var message model.Message

	result := repo.db.WithContext(ctx).First(&message, messageID)
	err := result.Error

	// 查询消息内容
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果消息不存在，返回业务错误
			return nil, utils.NewBusinessError(utils.ErrCodeResourceNotFound, "消息不存在或已被删除，请刷新页面后重试")
		}
		return nil, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return &message, nil
}

// GetUnreadMessageCount 获取未读消息数
func (repo *MessageRepositoryImpl) GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error) {
	var count int64

	// 构建查询
	query := repo.db.WithContext(ctx).
		Table("messages m").
		Joins("INNER JOIN user_message_mappings um ON um.message_id = m.id").
		Where("um.user_id = ?", userID).
		Where("um.is_read = ?", "N").   // 只统计未读消息
		Where("m.is_deleted = ?", "N"). // 只查询未删除的消息
		Where("um.is_deleted = ?", "N") // 只查询未删除的用户消息

	// 如果指定了消息类型，添加类型筛选
	if messageType != "" {
		query = query.Where("m.type = ?", messageType)
	}

	// 执行计数查询，使用distinct确保消息不被重复计数
	if err := query.Select("count(distinct um.id)").Count(&count).Error; err != nil {
		return 0, utils.NewSystemError(fmt.Errorf("查询未读消息数失败: %w", err))
	}

	return int(count), nil
}

// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
//func (repo *MessageRepositoryImpl) MarkAllMessagesAsRead(ctx context.Context, userID int) error {
//	result := repo.db.WithContext(ctx).
//		Model(&model.MessageUserMapping{}).
//		Where("user_id = ? AND is_read = ? AND is_deleted = ?", userID, "N", "N").
//		Updates(map[string]interface{}{
//			"is_read":   "Y",
//			"read_time": time.Now(),
//		})
//
//	if result.Error != nil {
//		return utils.NewSystemError(fmt.Errorf("一键已读操作失败: %w", result.Error))
//	}
//
//	return nil
//}

// ListMessageGroupsByUserID 查询用户消息群组列表
func (repo *MessageRepositoryImpl) ListMessageGroupsByUserID(ctx context.Context, page, pageSize int, userID int, typeCode string) ([]*dto.MessageGroupDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var results []*dto.MessageGroupDTO

	// 构建查询
	query := repo.db.WithContext(ctx).Table("user_message_groups umg").
		Select(`
				umg.id AS msg_group_id,
			  	umg.group_name,
			  	m.title AS latest_title,
			  	m.content AS latest_content,
			  	m.send_time AS latest_send_time,
				CASE
					WHEN umg.latest_msg_id > COALESCE(uum.last_read_msg_id, 0) THEN 'Y'
					ELSE 'N'
			  	END AS has_unread
		`).
		Joins("JOIN messages m ON m.id = umg.latest_msg_id").
		Joins("LEFT JOIN user_unread_marks uum ON uum.msg_group_id = umg.id AND uum.user_id = ?", userID).
		Where("umg.is_deleted = ?", "N"). // 只查询未删除的消息分组
		Where("m.is_deleted = ?", "N").   // 只查询未删除消息
		Order("m.send_time DESC")

	if typeCode == utils.TypeGroup {
		// 如果是群组消息类型，添加群组消息的筛选
		query = query.Joins("JOIN user_msg_group_mappings umgm ON umgm.msg_group_id = umg.id").
			Where("umgm.user_id = ?", userID).
			Where("umgm.is_deleted = ?", "N") // 只查询未退出的群组
	} else if typeCode == utils.TypeSystem {
		// 如果是系统消息类型，添加系统消息的筛选
		query = query.Where("umg.include_all_user = ?", "Y") // 查询包含所有用户的系统消息分组
	} else {
		// 否则返回参数不合法
		return nil, 0, utils.NewBusinessError(utils.ErrCodeParamInvalid, "消息类型参数不合法")
	}

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算总数时数据库查询失败: %v", err))
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return results, total, nil
}

// ListMsgByGroups 分页查询分组内消息列表
func (repo *MessageRepositoryImpl) ListMsgByGroups(ctx context.Context, page, pageSize int, groupID int) ([]*dto.ListMessageDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var results []*dto.ListMessageDTO
	query := repo.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("messages m").
		Select("m.id, m.title, m.content, m.send_time").
		Joins("JOIN message_group_mappings mgm ON mgm.message_id = m.id").
		Where("umgm.msg_group_id = ?", groupID).
		Where("m.is_deleted = ?", "N").  // 只查询未删除的消息
		Where("mgm.is_deleted = ?", "N") // 只查询未删除的组内消息

	// 按发送时间降序排列
	query = query.Order("m.send_time DESC")

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Select("count(distinct m.id)").Count(&total).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算总数时数据库查询失败: %v", err))
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return results, total, nil
}

// MarkAsReadByGroup 更新组内消息为已读
// 根据用户ID、群组ID、更新该用户在该群组内的最后已读消息ID
func (repo *MessageRepositoryImpl) MarkAsReadByGroup(ctx context.Context, userID int, groupID int) {
	// 获取组内最新消息ID
	latestMsgID, err := repo.getLatestMsgIDInGroup(ctx, groupID)
	if err != nil {
		// 只记录日志，更新已读状态失败不影响消息加载
		logrus.Errorf("获取组内最新消息ID失败: %v", err)
		return
	}
	// 更新用户在该组内的最后已读消息ID
	result := repo.db.WithContext(ctx).Model(&model.UserUnreadMark{}).
		Where("user_id = ? AND msg_group_id = ?", userID, groupID).
		Updates(map[string]interface{}{
			"last_read_msg_id": latestMsgID,
			"update_user":      userID,
		})
	if result.Error != nil {
		// 只记录日志，更新已读状态失败不影响消息加载
		logrus.Errorf("更新消息状态失败: %v", result.Error)
	}
}

// 获取组内最新消息ID
func (repo *MessageRepositoryImpl) getLatestMsgIDInGroup(ctx context.Context, groupID int) (int, error) {
	var latestMsgID int

	err := repo.db.WithContext(ctx).
		Table("message_group_mappings").
		Select("MAX(message_id)").
		Where("msg_group_id = ? AND is_deleted = ?", groupID, "N").
		Scan(&latestMsgID).Error

	if err != nil {
		return 0, utils.NewSystemError(fmt.Errorf("查询组内最新消息ID失败: %w", err))
	}

	return latestMsgID, nil
}

// CheckUserMsgPermission 权限校验查询，确保普通用户只能查看自己的消息
func (repo *MessageRepositoryImpl) CheckUserMsgPermission(ctx context.Context, userID int, groupID int) error {
	// 如果是全体用户消息组，直接返回true
	var count int64
	err := repo.db.WithContext(ctx).
		Table("user_message_groups").
		Where("id = ? AND include_all_user = ? AND is_deleted = ?", groupID, "Y", "N").
		Count(&count).Error
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("权限校验时数据库查询失败: %w", err))
	}

	if count == 1 {
		return nil
	}

	// 如果是管理员，直接返回true
	err = repo.db.WithContext(ctx).
		Table("users").
		Where("id = ? AND role = ?", userID, utils.RoleAdmin).
		Count(&count).Error
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("权限校验时数据库查询失败: %w", err))
	}

	if count == 1 {
		return nil
	}

	// 检查普通用户是否属于该消息组
	err = repo.db.WithContext(ctx).
		Table("user_msg_group_mappings").
		Where("user_id = ? AND msg_group_id = ? AND is_deleted = ?", userID, groupID, "N").
		Count(&count).Error

	if err != nil {
		return utils.NewSystemError(fmt.Errorf("权限校验时数据库查询失败: %w", err))
	}

	if count == 0 {
		return utils.NewBusinessError(utils.ErrCodePermissionDenied, "无权访问该消息组")
	}

	return nil
}
