package message

import (
	"context"
	msgmodel "news-release/internal/model/message"
	msgrepo "news-release/internal/repository/message"
)

// 服务接口，定义方法，接收 context.Context 和数据模型。
type MessageService interface {
	ListMessage(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*msgmodel.Message, int64, error)
	GetMessageContent(ctx context.Context, messageID int) (*msgmodel.Message, error)
}

// 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type MessageServiceImpl struct {
	messageRepo msgrepo.MessageRepository
}

// 创建服务实例
func NewMessageService(messageRepo msgrepo.MessageRepository) MessageService {
	return &MessageServiceImpl{messageRepo: messageRepo}
}

// 分页查询数据
func (s *MessageServiceImpl) ListMessage(ctx context.Context, page, pageSize int, userID int, messageType string) ([]*msgmodel.Message, int64, error) {
	return s.messageRepo.List(ctx, page, pageSize, userID, messageType)
}

// 获取消息内容
func (s *MessageServiceImpl) GetMessageContent(ctx context.Context, messageID int) (*msgmodel.Message, error) {
	return s.messageRepo.GetMessageContent(ctx, messageID)
}
