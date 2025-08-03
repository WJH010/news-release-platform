package service

import (
	"context"
	"news-release/internal/event/model"
	"news-release/internal/event/repository"
)

// EventService 定义事件服务接口，提供事件相关的业务逻辑方法
type EventService interface {
	ListEvent(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error)
}

// EventServiceImpl 实现 EventService 接口，提供事件相关的业务逻辑
type EventServiceImpl struct {
	eventRepo repository.EventRepository // 事件数据访问接口
}

// NewEventService 创建服务实例
func NewEventService(eventRepo repository.EventRepository) EventService {
	return &EventServiceImpl{eventRepo: eventRepo}
}

// ListEvent 分页查询事件列表
func (s *EventServiceImpl) ListEvent(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error) {
	return s.eventRepo.List(ctx, page, pageSize, eventStatus)
}
