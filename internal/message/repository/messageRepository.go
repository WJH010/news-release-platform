package repository

import (
	"context"
	"fmt"
	"news-release/internal/message/model"
	"time"

	"gorm.io/gorm"
)

// MessageRepository 数据访问接口，定义数据访问的方法集
type MessageRepository interface {
	// List 分页查询
	List(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*model.Message, int64, error)
	// GetMessageContent 内容查询
	GetMessageContent(ctx context.Context, messageID int) (*model.Message, error)
	// GetUnreadMessageCount 获取未读消息数
	GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error)
	// MarkAsRead 更新消息为已读
	MarkAsRead(ctx context.Context, userID, messageID int) error
	// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
	MarkAllMessagesAsRead(ctx context.Context, userID int) error
}

// MessageRepositoryImpl 实现接口的具体结构体
type MessageRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageRepository 创建数据访问实例
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

// List 分页查询数据
func (repo *MessageRepositoryImpl) List(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*model.Message, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var messages []*model.Message
	query := repo.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("messages m").
		Select("m.id, um.user_id, m.title, m.content, um.is_read, m.send_time, m.type, mt.type_name").
		Joins("INNER JOIN user_message_mappings um ON um.message_id = m.id").
		Joins("LEFT JOIN message_types mt ON m.type = mt.type_code").
		Where("u.user_id = ?", userID).
		Where("m.status = ?", 1)

	if messageType != "" {
		query = query.Where("m.type = ?", messageType)
	}

	// 按发送时间降序排列
	query = query.Order("m.send_time DESC")

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Select("count(distinct m.id)").Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算总数时数据库查询失败: %v", err)
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, fmt.Errorf("数据库查询失败: %v", err)
	}

	return messages, total, nil
}

// GetMessageContent 内容查询
func (repo *MessageRepositoryImpl) GetMessageContent(ctx context.Context, messageID int) (*model.Message, error) {
	var message model.Message

	result := repo.db.WithContext(ctx).First(&message, messageID)
	err := result.Error

	// 查询消息内容
	if err != nil {
		return nil, err
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
		Where("um.is_read = ?", "N"). // 只统计未读消息
		Where("m.status = ?", 1)

	// 如果指定了消息类型，添加类型筛选
	if messageType != "" {
		query = query.Where("m.type = ?", messageType)
	}

	// 执行计数查询，使用distinct确保消息不被重复计数
	if err := query.Select("count(distinct um.id)").Count(&count).Error; err != nil {
		return 0, fmt.Errorf("查询未读消息数失败: %v", err)
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
		return fmt.Errorf("更新消息状态失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
func (repo *MessageRepositoryImpl) MarkAllMessagesAsRead(ctx context.Context, userID int) error {
	result := repo.db.WithContext(ctx).
		Model(&model.MessageUserMapping{}).
		Where("user_id = ? AND is_read = ? AND status = ?", userID, "N", 1).
		Updates(map[string]interface{}{
			"is_read":   "Y",
			"read_time": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("一键已读操作失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
