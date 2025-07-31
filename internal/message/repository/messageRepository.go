package repository

import (
	"context"
	"fmt"
	"news-release/internal/message/model"
	"time"

	"gorm.io/gorm"
)

// 数据访问接口，定义数据访问的方法集
type MessageRepository interface {
	// 分页查询
	List(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*model.Message, int64, error)
	// 内容查询
	GetMessageContent(ctx context.Context, messageID int) (*model.Message, error)
	// 获取未读消息数
	GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error)
	// 更新消息为已读
	MarkAsRead(ctx context.Context, userID, messageID int) error
}

// 实现接口的具体结构体
type MessageRepositoryImpl struct {
	db *gorm.DB
}

// 创建数据访问实例
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

// 分页查询数据
func (r *MessageRepositoryImpl) List(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*model.Message, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var messages []*model.Message
	query := r.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("users u").
		Select("m.id, u.user_id, m.title, m.content, um.is_read, m.send_time, m.type, mt.type_name").
		Joins("INNER JOIN user_message_mappings um ON u.user_id = um.user_id").
		Joins("INNER JOIN messages m ON um.message_id = m.id").
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

// 内容查询
func (r *MessageRepositoryImpl) GetMessageContent(ctx context.Context, messageID int) (*model.Message, error) {
	var message model.Message

	result := r.db.WithContext(ctx).First(&message, messageID)
	err := result.Error

	// 查询消息内容
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// 获取未读消息数
func (r *MessageRepositoryImpl) GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error) {
	var count int64

	// 构建查询
	query := r.db.WithContext(ctx).
		Table("users u").
		Joins("INNER JOIN user_message_mappings um ON u.user_id = um.user_id").
		Joins("INNER JOIN messages m ON um.message_id = m.id").
		Where("u.user_id = ?", userID).
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

// 更新消息为已读
func (r *MessageRepositoryImpl) MarkAsRead(ctx context.Context, userID, messageID int) error {
	result := r.db.WithContext(ctx).
		Model(&model.UserMessageMapping{}).
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
