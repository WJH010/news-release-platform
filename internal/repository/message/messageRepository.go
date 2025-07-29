package message

import (
	"context"
	"fmt"
	msgmodel "news-release/internal/model/message"

	"gorm.io/gorm"
)

// 数据访问接口，定义数据访问的方法集
type MessageRepository interface {
	// 分页查询
	List(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*msgmodel.Message, int64, error)
	// 内容查询
	GetMessageContent(ctx context.Context, messageID int) (*msgmodel.Message, error)
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
func (r *MessageRepositoryImpl) List(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*msgmodel.Message, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var messages []*msgmodel.Message
	query := r.db.WithContext(ctx)

	// 构建基础查询
	query = query.Table("users u").
		Select("u.user_id, m.title, m.content, um.is_read, m.send_time, mt.type_name").
		Joins("INNER JOIN user_message_mappings um ON u.user_id = um.user_id").
		Joins("INNER JOIN messages m ON um.message_id = m.id").
		Joins("INNER JOIN message_types mt ON m.type = mt.type_code").
		Where("u.user_id = ?", userID)

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
func (r *MessageRepositoryImpl) GetMessageContent(ctx context.Context, messageID int) (*msgmodel.Message, error) {
	var message msgmodel.Message

	result := r.db.WithContext(ctx).First(&message, messageID)
	err := result.Error

	// 查询消息内容
	if err != nil {
		return nil, err
	}

	return &message, nil
}
