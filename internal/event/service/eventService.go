package service

import (
	"context"
	"news-release/internal/event/model"
	"news-release/internal/event/repository"
	userrepo "news-release/internal/user/repository"
	"news-release/internal/utils"
	"time"
)

// EventService 定义事件服务接口，提供事件相关的业务逻辑方法
type EventService interface {
	// ListEvent 分页查询活动列表
	ListEvent(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error)
	// GetEventDetail 获取活动详情
	GetEventDetail(ctx context.Context, eventID int) (*model.Event, error)
	// RegistrationEvent 活动报名
	RegistrationEvent(ctx context.Context, eventID int, userID int) error
	// IsUserRegistered 查询用户是否已报名活动
	IsUserRegistered(ctx context.Context, eventID int, userID int) (bool, error)
}

// EventServiceImpl 实现 EventService 接口，提供事件相关的业务逻辑
type EventServiceImpl struct {
	eventRepo repository.EventRepository // 事件数据访问接口
	userRepo  userrepo.UserRepository    // 用户数据访问接口
}

// NewEventService 创建服务实例
func NewEventService(
	eventRepo repository.EventRepository,
	userRepo userrepo.UserRepository,
) EventService {
	return &EventServiceImpl{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

// ListEvent 分页查询活动列表
func (svc *EventServiceImpl) ListEvent(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error) {
	return svc.eventRepo.List(ctx, page, pageSize, eventStatus)
}

// GetEventDetail 获取活动详情
func (svc *EventServiceImpl) GetEventDetail(ctx context.Context, eventID int) (*model.Event, error) {
	event, err := svc.eventRepo.GetEventDetail(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// 获取关联图片列表
	var images []repository.EventImage
	images = svc.eventRepo.ListEventImage(ctx, eventID)

	// 添加图片到活动详情
	event.Images = make([]string, 0, len(images)) // 预分配空间，提高性能
	for _, img := range images {
		event.Images = append(event.Images, img.URL)
	}

	return event, nil
}

// RegistrationEvent 活动报名实现
func (svc *EventServiceImpl) RegistrationEvent(ctx context.Context, eventID int, userID int) error {
	var mapping *model.EventUserMapping
	// 检查活动是否存在
	event, err := svc.eventRepo.GetEventDetail(ctx, eventID)
	if err != nil {
		return err
	}
	// 检查活动是否已删除
	if event.IsDeleted == utils.DeletedFlagYes {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "活动已失效")
	}
	// 检查活动是否在报名时间内
	if event.RegistrationStartTime.After(time.Now()) || event.RegistrationEndTime.Before(time.Now()) {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "未在活动报名时间内")
	}
	// 检查用户信息是否完整
	user, err := svc.userRepo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "加载用户信息失败")
	}
	if user.Name == "" || user.PhoneNumber == "" || user.Email == "" || user.Unit == "" || user.Department == "" || user.Position == "" || user.Industry == "" {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "用户信息不完整，请完善个人信息")
	}

	// 执行活动报名逻辑
	mapping, err = svc.eventRepo.GetEventUserMap(ctx, eventID, userID)
	if err != nil {
		return err
	}
	// 如果关联关系不存在，则创建新的关联关系
	if mapping == nil {
		mapping = &model.EventUserMapping{
			UserID:  userID,
			EventID: eventID,
		}
		err = svc.eventRepo.CreatEventUserMap(ctx, mapping)
		if err != nil {
			return err
		}
	} else if mapping.IsDeleted == utils.DeletedFlagYes {
		// 如果关联关系软删除了，则恢复
		err = svc.eventRepo.UpdateEUMapDeleteFlag(ctx, eventID, userID, utils.DeletedFlagNo)
		if err != nil {
			return err
		}
	} else if mapping.IsDeleted == utils.DeletedFlagNo {
		// 如果关联关系存在且有效，则返回错误提示
		return utils.NewBusinessError(utils.ErrCodeResourceExists, "已报名该活动，请勿重复报名")
	}

	return nil
}

// IsUserRegistered 查询用户是否已报名活动
func (svc *EventServiceImpl) IsUserRegistered(ctx context.Context, eventID int, userID int) (bool, error) {
	return svc.eventRepo.IsUserRegistered(ctx, eventID, userID)
}
