package service

import (
	"context"
	"news-release/internal/message/dto"
	"news-release/internal/message/model"
	"news-release/internal/message/repository"
)

// MessageService 服务接口，定义方法，接收 context.Context 和数据模型。
type MessageService interface {
	// GetMessageContent 获取消息内容
	GetMessageContent(ctx context.Context, messageID int) (*model.Message, error)
	// GetUnreadMessageCount 获取未读消息数
	GetUnreadMessageCount(ctx context.Context, userID int, messageType string) (int, error)
	// MarkAsRead 标记消息为已读
	//MarkAsRead(ctx context.Context, userID, messageID int) error
	// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
	//MarkAllMessagesAsRead(ctx context.Context, userID int) error
	// ListMessageGroupsByUserID 分页查询用户消息群组列表
	ListMessageGroupsByUserID(ctx context.Context, page, pageSize int, userID int, typeCode string) ([]*dto.MessageGroupDTO, int64, error)
	// ListMsgByGroups 分页查询分组内消息列表
	ListMsgByGroups(ctx context.Context, page, pageSize int, userID int, groupID int) ([]*dto.ListMessageDTO, int64, error)
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
//func (svc *MessageServiceImpl) MarkAsRead(ctx context.Context, userID, messageID int) error {
//	return svc.messageRepo.MarkAsRead(ctx, userID, messageID)
//}

// MarkAllMessagesAsRead 一键已读，更新所有未读消息为已读
//func (svc *MessageServiceImpl) MarkAllMessagesAsRead(ctx context.Context, userID int) error {
//	return svc.messageRepo.MarkAllMessagesAsRead(ctx, userID)
//}

// ListMessageGroupsByUserID 分页查询用户消息群组列表
func (svc *MessageServiceImpl) ListMessageGroupsByUserID(ctx context.Context, page, pageSize int, userID int, typeCode string) ([]*dto.MessageGroupDTO, int64, error) {
	return svc.messageRepo.ListMessageGroupsByUserID(ctx, page, pageSize, userID, typeCode)
}

// ListMsgByGroups 分页查询分组内消息列表
func (svc *MessageServiceImpl) ListMsgByGroups(ctx context.Context, page, pageSize int, userID int, groupID int) ([]*dto.ListMessageDTO, int64, error) {
	// 校验权限，确保普通用户只能查看自己的消息
	err := svc.messageRepo.CheckUserMsgPermission(ctx, userID, groupID)
	if err != nil {
		return nil, 0, err
	}
	// 标记组内所有消息为已读
	svc.messageRepo.MarkAsReadByGroup(ctx, userID, groupID)

	return svc.messageRepo.ListMsgByGroups(ctx, page, pageSize, groupID)
}
