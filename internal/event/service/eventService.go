package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	db "news-release/internal/database"
	"news-release/internal/event/dto"
	"news-release/internal/event/model"
	"news-release/internal/event/repository"
	filerepo "news-release/internal/file/repository"
	usermodel "news-release/internal/user/model"
	userrepo "news-release/internal/user/repository"
	"news-release/internal/utils"
	"time"
)

// EventService 定义事件服务接口，提供事件相关的业务逻辑方法
type EventService interface {
	// ListEvent 分页查询活动列表
	ListEvent(ctx context.Context, page, pageSize int, eventStatus string, isDeleted string) ([]*model.Event, int, error)
	// GetEventDetail 获取活动详情
	GetEventDetail(ctx context.Context, eventID int) (*model.Event, error)
	// RegistrationEvent 活动报名
	RegistrationEvent(ctx context.Context, eventID int, userID int) error
	// CancelRegistrationEvent 取消活动报名
	CancelRegistrationEvent(ctx context.Context, eventID int, userID int) error
	// IsUserRegistered 查询用户是否已报名活动
	IsUserRegistered(ctx context.Context, eventID int, userID int) (bool, error)
	// ListUserRegisteredEvents 获取用户已报名的活动列表
	ListUserRegisteredEvents(ctx context.Context, page, pageSize int, userID int, eventStatus string) ([]*model.Event, int, error)
	// CreateEvent 创建活动
	CreateEvent(ctx context.Context, event *model.Event, imageIDList []int) error
	// UpdateEvent 更新活动
	UpdateEvent(ctx context.Context, eventID int, req dto.UpdateEventRequest, userID int) error
	// DeleteEvent 删除活动
	DeleteEvent(ctx context.Context, eventID int, userID int) error
	// ListEventRegisteredUser 获取活动报名用户列表
	ListEventRegisteredUser(ctx context.Context, page, pageSize int, eventID int) ([]*usermodel.User, int, error)
}

// EventServiceImpl 实现 EventService 接口，提供事件相关的业务逻辑
type EventServiceImpl struct {
	eventRepo repository.EventRepository // 事件数据访问接口
	userRepo  userrepo.UserRepository    // 用户数据访问接口
	fileRepo  filerepo.FileRepository    // 文件数据访问接口
}

// NewEventService 创建服务实例
func NewEventService(
	eventRepo repository.EventRepository,
	userRepo userrepo.UserRepository,
	fileRepo filerepo.FileRepository,
) EventService {
	return &EventServiceImpl{
		eventRepo: eventRepo,
		userRepo:  userRepo,
		fileRepo:  fileRepo,
	}
}

// ListEvent 分页查询活动列表
func (svc *EventServiceImpl) ListEvent(ctx context.Context, page, pageSize int, eventStatus string, isDeleted string) ([]*model.Event, int, error) {
	return svc.eventRepo.List(ctx, page, pageSize, eventStatus, isDeleted)
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

// CancelRegistrationEvent 取消活动报名
func (svc *EventServiceImpl) CancelRegistrationEvent(ctx context.Context, eventID int, userID int) error {
	// 检查活动是否存在
	event, err := svc.eventRepo.GetEventDetail(ctx, eventID)
	if err != nil {
		return err
	}
	// 检查活动是否已删除
	if event.IsDeleted == utils.DeletedFlagYes {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "活动已失效")
	}

	// 执行取消报名逻辑
	err = svc.eventRepo.UpdateEUMapDeleteFlag(ctx, eventID, userID, utils.DeletedFlagYes)
	if err != nil {
		return err
	}

	return nil
}

// IsUserRegistered 查询用户是否已报名活动
func (svc *EventServiceImpl) IsUserRegistered(ctx context.Context, eventID int, userID int) (bool, error) {
	return svc.eventRepo.IsUserRegistered(ctx, eventID, userID)
}

// ListUserRegisteredEvents 获取用户已报名的活动列表
func (svc *EventServiceImpl) ListUserRegisteredEvents(ctx context.Context, page, pageSize int, userID int, eventStatus string) ([]*model.Event, int, error) {
	return svc.eventRepo.ListUserRegisteredEvents(ctx, page, pageSize, userID, eventStatus)
}

// CreateEvent 创建活动
func (svc *EventServiceImpl) CreateEvent(ctx context.Context, event *model.Event, imageIDList []int) error {
	// 检查活动时间是否合理
	if event.EventStartTime.After(event.EventEndTime) {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "活动开始时间不能晚于结束时间")
	}
	if event.RegistrationStartTime.After(event.RegistrationEndTime) {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "报名开始时间不能晚于结束时间")
	}

	// 开启事务
	tx := db.GetDB().Begin()
	if tx.Error != nil {
		return utils.NewSystemError(fmt.Errorf("开启事务失败: %w", tx.Error))
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logrus.Panic("事务回滚，发生异常: ", r)
		}
	}()

	// 创建活动
	if err := svc.eventRepo.CreateEvent(ctx, tx, event); err != nil {
		tx.Rollback()
		return err
	}

	// 如果有图片，更新images表的biz_id和biz_type
	if len(imageIDList) > 0 {
		if err := svc.fileRepo.BatchUpdateImageBizID(ctx, tx, imageIDList, event.ID, utils.TypeEvent); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.NewSystemError(fmt.Errorf("提交事务失败: %w", err))
	}
	return nil
}

// UpdateEvent 更新活动
func (svc *EventServiceImpl) UpdateEvent(ctx context.Context, eventID int, req dto.UpdateEventRequest, userID int) error {
	// 检查活动是否存在
	event, err := svc.eventRepo.GetEventDetail(ctx, eventID)
	if err != nil {
		return err
	}

	// 构建更新字段映射
	updateFields, err := makeUpdateFields(event, req)
	if err != nil {
		return err
	}

	var imageIDList []int
	if req.ImageIDList != nil {
		imageIDList = *req.ImageIDList
	}

	// 设置更新人
	updateFields["update_user"] = userID

	// 开启事务
	tx := db.GetDB().Begin()
	if tx.Error != nil {
		return utils.NewSystemError(fmt.Errorf("开启事务失败: %w", tx.Error))
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logrus.Panic("事务回滚，发生异常: ", r)
		}
	}()

	// 更新活动
	if err := svc.eventRepo.UpdateEvent(ctx, tx, eventID, updateFields); err != nil {
		tx.Rollback()
		return err
	}

	// 如果有图片，更新images表的biz_id和biz_type
	if len(imageIDList) > 0 {
		if err := svc.fileRepo.BatchUpdateImageBizID(ctx, tx, imageIDList, eventID, utils.TypeEvent); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.NewSystemError(fmt.Errorf("提交事务失败: %w", err))
	}
	return nil
}

// makeUpdateFields 构建更新字段映射
func makeUpdateFields(event *model.Event, req dto.UpdateEventRequest) (map[string]interface{}, error) {
	updateFields := make(map[string]interface{})

	// 先处理时间字段，然后校验时间是否合理
	var eventStartTime, eventEndTime, registrationStartTime, registrationEndTime time.Time
	var err error
	if req.EventStartTime != nil {
		eventStartTime, err = utils.StringToTime(*req.EventStartTime)
		if err != nil {
			return nil, err
		}
		updateFields["event_start_time"] = eventStartTime
	} else {
		eventStartTime = event.EventStartTime
	}
	if req.EventEndTime != nil {
		eventEndTime, err = utils.StringToTime(*req.EventEndTime)
		if err != nil {
			return nil, err
		}
		updateFields["event_end_time"] = eventEndTime
	} else {
		eventEndTime = event.EventEndTime
	}
	// 检查活动时间是否合理
	if eventStartTime.After(eventEndTime) {
		return nil, utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "活动开始时间不能晚于结束时间")
	}
	if req.RegistrationStartTime != nil {
		registrationStartTime, err = utils.StringToTime(*req.RegistrationStartTime)
		if err != nil {
			return nil, err
		}
		updateFields["registration_start_time"] = registrationStartTime

	} else {
		registrationStartTime = event.RegistrationStartTime
	}
	if req.RegistrationEndTime != nil {
		registrationEndTime, err = utils.StringToTime(*req.RegistrationEndTime)
		if err != nil {
			return nil, err
		}
		updateFields["registration_end_time"] = registrationEndTime
	} else {
		registrationEndTime = event.RegistrationEndTime
	}
	// 检查报名时间是否合理
	if registrationStartTime.After(registrationEndTime) {
		return nil, utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "报名开始时间不能晚于结束时间")
	}

	if req.Title != nil {
		updateFields["title"] = *req.Title
	}
	if req.Detail != nil {
		updateFields["detail"] = *req.Detail
	}
	if req.EventAddress != nil {
		updateFields["event_address"] = *req.EventAddress
	}
	if req.RegistrationFee != nil {
		updateFields["registration_fee"] = *req.RegistrationFee
	}
	if req.CoverImageURL != nil {
		updateFields["cover_image_url"] = *req.CoverImageURL
	}

	return updateFields, nil
}

// DeleteEvent 删除活动
func (svc *EventServiceImpl) DeleteEvent(ctx context.Context, eventID int, userID int) error {
	// 检查活动是否存在
	event, err := svc.eventRepo.GetEventDetail(ctx, eventID)
	if err != nil {
		return err
	}
	// 检查活动是否已删除
	if event.IsDeleted == utils.DeletedFlagYes {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "活动已失效")
	}

	// 执行软删除逻辑
	err = svc.eventRepo.DeleteEvent(ctx, eventID, userID)
	if err != nil {
		return err
	}

	return nil
}

// ListEventRegisteredUser 获取活动报名用户列表
func (svc *EventServiceImpl) ListEventRegisteredUser(ctx context.Context, page, pageSize int, eventID int) ([]*usermodel.User, int, error) {
	return svc.eventRepo.ListEventRegisteredUser(ctx, page, pageSize, eventID)
}
