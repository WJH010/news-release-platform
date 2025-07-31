package service

import (
	"context"
	"news-release/internal/message/model"
	"news-release/internal/message/repository"
)

// 服务接口，定义方法，接收 context.Context 和数据模型。
type MessageService interface {
	ListMessage(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*model.Message, int64, error)
	GetMessageContent(ctx context.Context, messageID int) (*model.Message, error)
	GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error)
}

// 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type MessageServiceImpl struct {
	messageRepo repository.MessageRepository
}

// 创建服务实例
func NewMessageService(messageRepo repository.MessageRepository) MessageService {
	return &MessageServiceImpl{messageRepo: messageRepo}
}

// 分页查询数据
func (s *MessageServiceImpl) ListMessage(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*model.Message, int64, error) {
	return s.messageRepo.List(ctx, page, pageSize, userID, messageType)
}

// 获取消息内容
func (s *MessageServiceImpl) GetMessageContent(ctx context.Context, messageID int) (*model.Message, error) {
	return s.messageRepo.GetMessageContent(ctx, messageID)
}

// 获取未读消息数
func (s *MessageServiceImpl) GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error) {
	return s.messageRepo.GetUnreadMessageCount(ctx, userID, messageType)
}
