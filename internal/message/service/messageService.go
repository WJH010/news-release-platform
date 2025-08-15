package service

import (
	"context"
	"news-release/internal/message/dto"
	"news-release/internal/message/model"
	"news-release/internal/message/repository"
	"news-release/internal/utils"
)

// MessageService 服务接口，定义方法，接收 context.Context 和数据模型。
type MessageService interface {
	GetMessageContent(ctx context.Context, messageID int) (*model.Message, error)
	GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error)
	MarkAsRead(ctx context.Context, userID, messageID int) error
	MarkAllMessagesAsRead(ctx context.Context, userID int) error
	ListMessageByTypeGroups(ctx context.Context, page, pageSize int, userID int, typeCodes []string) ([]*dto.MessageGroupDTO, int64, error)
	ListMessageByEventGroups(ctx context.Context, page, pageSize int, userID int) ([]*dto.MessageGroupDTO, int64, error)
	ListMsgByGroups(ctx context.Context, page, pageSize int, userID int, eventID int, messageType string) ([]*dto.ListMessageDTO, int64, error)
}

// MessageServiceImpl 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type MessageServiceImpl struct {
	messageRepo repository.MessageRepository
}

// NewMessageService 创建服务实例
func NewMessageService(messageRepo repository.MessageRepository) MessageService {
	return &MessageServiceImpl{messageRepo: messageRepo}
}

// GetMessageContent 获取消息内容
func (svc *MessageServiceImpl) GetMessageContent(ctx context.Context, messageID int) (*model.Message, error) {
	return svc.messageRepo.GetMessageContent(ctx, messageID)
}

// GetUnreadMessageCount 获取未读消息数
func (svc *MessageServiceImpl) GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error) {
	return svc.messageRepo.GetUnreadMessageCount(ctx, userID, messageType)
}

// MarkAsRead 标记消息为已读
func (svc *MessageServiceImpl) MarkAsRead(ctx context.Context, userID, messageID int) error {
	return svc.messageRepo.MarkAsRead(ctx, userID, messageID)
}

// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
func (svc *MessageServiceImpl) MarkAllMessagesAsRead(ctx context.Context, userID int) error {
	return svc.messageRepo.MarkAllMessagesAsRead(ctx, userID)
}

// ListMessageByTypeGroups 按类型分组查询消息
func (svc *MessageServiceImpl) ListMessageByTypeGroups(ctx context.Context, page, pageSize int, userID int, typeCodes []string) ([]*dto.MessageGroupDTO, int64, error) {
	return svc.messageRepo.ListMessageByTypeGroups(ctx, page, pageSize, userID, typeCodes)
}

// ListMessageByEventGroups 按活动分组查询消息
func (svc *MessageServiceImpl) ListMessageByEventGroups(ctx context.Context, page, pageSize int, userID int) ([]*dto.MessageGroupDTO, int64, error) {
	return svc.messageRepo.ListMessageByEventGroups(ctx, page, pageSize, userID)
}

// ListMsgByGroups 分页查询分组内消息列表
func (svc *MessageServiceImpl) ListMsgByGroups(ctx context.Context, page, pageSize int, userID int, eventID int, messageType string) ([]*dto.ListMessageDTO, int64, error) {
	if messageType == model.MessageTypeEvent && eventID < 1 {
		// 如果是活动消息类型，但未指定活动ID，返回错误
		return nil, 0, utils.NewBusinessError(utils.ErrCodeParamInvalid, "活动消息类型必须指定活动ID")
	}
	return svc.messageRepo.ListMsgByGroups(ctx, page, pageSize, userID, eventID, messageType)
}
