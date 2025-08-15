package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"news-release/internal/message/dto"
	"news-release/internal/message/model"
	"news-release/internal/utils"
	"time"

	"gorm.io/gorm"
)

// MessageRepository 数据访问接口，定义数据访问的方法集
type MessageRepository interface {
	// GetMessageContent 内容查询
	GetMessageContent(ctx context.Context, messageID int) (*model.Message, error)
	// GetUnreadMessageCount 获取未读消息数
	GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error)
	// MarkAsRead 更新消息为已读
	MarkAsRead(ctx context.Context, userID, messageID int) error
	// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
	MarkAllMessagesAsRead(ctx context.Context, userID int) error
	// ListMessageByTypeGroups 按消息类型分组统计消息列表
	ListMessageByTypeGroups(ctx context.Context, page, pageSize int, userID int, typeCodes []string) ([]*dto.MessageGroupDTO, int64, error)
	// ListMessageByEventGroups 按活动分组统计消息列表
	ListMessageByEventGroups(ctx context.Context, page, pageSize int, userID int) ([]*dto.MessageGroupDTO, int64, error)
	// ListMsgByGroups 分页查询分组内消息列表
	ListMsgByGroups(ctx context.Context, page, pageSize int, userID int, eventID int, messageType string) ([]*dto.ListMessageDTO, int64, error)
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

// MarkAsRead 更新消息为已读
func (repo *MessageRepositoryImpl) MarkAsRead(ctx context.Context, userID, messageID int) error {
	result := repo.db.WithContext(ctx).
		Model(&model.MessageUserMapping{}).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Updates(map[string]interface{}{
			"is_read":     "Y",
			"read_time":   time.Now(),
			"update_time": time.Now(),
		})

	if result.Error != nil {
		// return utils.NewSystemError(fmt.Errorf("更新消息状态失败: %w", result.Error))
		// 只记录日志，更新已读状态失败不影响消息加载
		logrus.Errorf("更新消息状态失败: %v", result.Error)
	}
	return nil
}

// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
func (repo *MessageRepositoryImpl) MarkAllMessagesAsRead(ctx context.Context, userID int) error {
	result := repo.db.WithContext(ctx).
		Model(&model.MessageUserMapping{}).
		Where("user_id = ? AND is_read = ? AND is_deleted = ?", userID, "N", "N").
		Updates(map[string]interface{}{
			"is_read":   "Y",
			"read_time": time.Now(),
		})

	if result.Error != nil {
		return utils.NewSystemError(fmt.Errorf("一键已读操作失败: %w", result.Error))
	}

	return nil
}

// ListMessageByTypeGroups 按消息类型分组统计消息列表
// 该方法查询每个消息类型的最新一条消息，并统计未读消息数
func (repo *MessageRepositoryImpl) ListMessageByTypeGroups(ctx context.Context, page, pageSize int, userID int, typeCodes []string) ([]*dto.MessageGroupDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var results []*dto.MessageGroupDTO

	// 构建子查询
	subQuery := repo.db.Table("messages m").
		Select(`
			mt.type_name AS group_name,
			COUNT(CASE WHEN mum.is_read = 'N' THEN 1 END) OVER (PARTITION BY mt.id) AS unread_count,
			m.content AS latest_content,
			m.send_time AS latest_time,
			ROW_NUMBER() OVER (PARTITION BY mt.id ORDER BY m.send_time DESC) AS rn
		`).
		Joins("JOIN message_user_mappings mum ON m.id = mum.message_id").
		Joins("JOIN message_types mt ON m.type = mt.type_code").
		Where("mum.user_id = ?", userID).
		Where("m.type IN (?)", typeCodes).
		Where("m.is_deleted = ?", "N").
		Where("mum.is_deleted = ?", "N")

	// 主查询：筛选每组最新的一条消息
	query := repo.db.WithContext(ctx).
		Table("(?) AS t", subQuery).
		Where("t.rn = 1").
		Select("t.group_name, t.unread_count, t.latest_content, t.latest_time").
		Order("t.latest_time DESC")

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

// ListMessageByEventGroups 按活动分组统计消息列表
// 该方法查询每个活动的最新一条消息，并统计未读消息数
func (repo *MessageRepositoryImpl) ListMessageByEventGroups(ctx context.Context, page, pageSize int, userID int) ([]*dto.MessageGroupDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var results []*dto.MessageGroupDTO

	// 构建子查询
	subQuery := repo.db.Table("messages m").
		Select(`
			e.title AS group_name,
			COUNT(CASE WHEN mum.is_read = 'N' THEN 1 END) OVER (PARTITION BY e.id) AS unread_count,
			m.content AS latest_content,
			m.send_time AS latest_time,
			ROW_NUMBER() OVER (PARTITION BY e.id ORDER BY m.send_time DESC) AS rn
		`).
		Joins("JOIN message_user_mappings mum ON m.id = mum.message_id").
		Joins("JOIN message_event_mappings mem ON m.id = mem.message_id").
		Joins("JOIN events e ON mem.event_id = e.id").
		Where("mum.user_id = ?", userID).
		Where("m.type = ?", "EVENT").
		Where("m.is_deleted = ?", "N").
		Where("mum.is_deleted = ?", "N").
		Where("e.is_deleted = ?", "N").
		Where("mem.is_deleted = ?", "N")

	// 主查询：筛选每组最新的一条消息
	query := repo.db.WithContext(ctx).
		Table("(?) AS t", subQuery).
		Where("t.rn = 1").
		Select("t.group_name, t.unread_count, t.latest_content, t.latest_time").
		Order("t.latest_time DESC")

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
func (repo *MessageRepositoryImpl) ListMsgByGroups(ctx context.Context, page, pageSize int, userID int, eventID int, messageType string) ([]*dto.ListMessageDTO, int64, error) {
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
		Select("m.id, m.title, m.content, mum.is_read, m.send_time").
		Joins("JOIN message_user_mappings mum ON mum.message_id = m.id").
		Where("mum.user_id = ?", userID).
		Where("m.type = ?", messageType).
		Where("m.is_deleted = ?", "N").  // 只查询未删除的消息
		Where("mum.is_deleted = ?", "N") // 只查询未删除的用户消息

	// 如果是活动消息，添加活动ID筛选
	if messageType == model.MessageTypeEvent && eventID > 0 {
		query = query.Joins("JOIN message_event_mappings mem ON mem.message_id = m.id").
			Where("mem.event_id = ?", eventID).
			Where("mem.is_deleted = ?", "N") // 只查询未删除的活动消息映射
	}

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
