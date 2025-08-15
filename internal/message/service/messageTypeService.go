package service

import (
	"context"
	"news-release/internal/message/model"
	"news-release/internal/message/repository"
)

// MessageTypeService 定义消息类型服务接口，提供消息类型相关的业务逻辑方法
type MessageTypeService interface {
	// ListMessageType 列出所有消息类型
	ListMessageType(ctx context.Context) ([]*model.MessageType, error)
}

// MessageTypeServiceImpl 实现 messageTypeService 接口
type MessageTypeServiceImpl struct {
	messageTypeRepo repository.MessageTypeRepository // 消息类型数据访问接口
}

// NewMessageTypeService 创建消息类型服务实例
func NewMessageTypeService(messageTypeRepo repository.MessageTypeRepository) MessageTypeService {
	return &MessageTypeServiceImpl{
		messageTypeRepo: messageTypeRepo,
	}
}

// ListMessageType 列出所有消息类型
func (svc *MessageTypeServiceImpl) ListMessageType(ctx context.Context) ([]*model.MessageType, error) {
	return svc.messageTypeRepo.List(ctx)
}
