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

// ListEvent 分页查询活动列表
func (s *EventServiceImpl) ListEvent(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error) {
	events, total, err := s.eventRepo.List(ctx, page, pageSize, eventStatus)
	if err != nil {
		return nil, 0, err
	}

	// 获取关联图片列表
	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}
	var images []repository.EventImage
	if len(eventIDs) > 0 {
		images = s.eventRepo.ListEventImage(ctx, eventIDs)
	}

	// 建立图片映射关系
	imageMap := make(map[int][]string)
	for _, img := range images {
		imageMap[img.BizID] = append(imageMap[img.BizID], img.URL)
	}

	// 为每个活动设置图片列表
	for _, event := range events {
		event.Images = imageMap[event.ID]
	}

	return events, total, nil
}
