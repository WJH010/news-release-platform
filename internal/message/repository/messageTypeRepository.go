package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"news-release/internal/message/model"
	"news-release/internal/utils"
)

type MessageTypeRepository interface {
	// List 列出所有消息类型
	List(ctx context.Context) ([]*model.MessageType, error)
}

// MessageTypeRepositoryImpl 实现 MessageTypeRepository 接口
type MessageTypeRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageTypeRepository 创建消息类型数据访问实例
func NewMessageTypeRepository(db *gorm.DB) MessageTypeRepository {
	return &MessageTypeRepositoryImpl{db: db}
}

// List 列出所有消息类型
func (repo *MessageTypeRepositoryImpl) List(ctx context.Context) ([]*model.MessageType, error) {
	var messageTypes []*model.MessageType
	// 查询消息类型列表
	result := repo.db.WithContext(ctx).Find(&messageTypes)
	if result.Error != nil {
		return nil, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", result.Error))
	}

	return messageTypes, nil // 返回查询结果
}
