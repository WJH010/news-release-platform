package service

import (
	"context"
	"news-release/internal/event/model"
	"news-release/internal/event/repository"
)

// EventService 定义事件服务接口，提供事件相关的业务逻辑方法
type EventService interface {
	ListEvent(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error)
	GetEventDetail(ctx context.Context, eventID int) (*model.Event, error)
}

// EventServiceImpl 实现 EventService 接口，提供事件相关的业务逻辑
type EventServiceImpl struct {
	eventRepo repository.EventRepository // 事件数据访问接口
}

// NewEventService 创建服务实例
func NewEventService(eventRepo repository.EventRepository) EventService {
	return &EventServiceImpl{eventRepo: eventRepo}
}

// ListEvent 分页查询活动列表
func (s *EventServiceImpl) ListEvent(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error) {
	return s.eventRepo.List(ctx, page, pageSize, eventStatus)
}

// GetEventDetail 获取活动详情
func (s *EventServiceImpl) GetEventDetail(ctx context.Context, eventID int) (*model.Event, error) {
	event, err := s.eventRepo.GetEventDetail(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// 获取关联图片列表
	var images []repository.EventImage
	images = s.eventRepo.ListEventImage(ctx, eventID)

	// 添加图片到活动详情
	event.Images = make([]string, 0, len(images)) // 预分配空间，提高性能
	for _, img := range images {
		event.Images = append(event.Images, img.URL)
	}

	return event, nil
}
