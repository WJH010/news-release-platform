package service

import (
	"context"
	"fmt"
	"news-release/internal/event/dto"
	"news-release/internal/event/model"
	"news-release/internal/event/repository"
	filerepo "news-release/internal/file/repository"
	msgmodel "news-release/internal/message/model"
	msgsvc "news-release/internal/message/service"
	userrepo "news-release/internal/user/repository"
	"news-release/internal/utils"
	"time"

	"gorm.io/gorm"
)

// EventService 定义事件服务接口，提供事件相关的业务逻辑方法
type EventService interface {
	// GetEventStatus 根据开始时间和结束时间计算活动状态
	GetEventStatus(registrationStartTime time.Time, registrationEndTime time.Time) string
	// ListEvent 分页查询活动列表
	ListEvent(ctx context.Context, page, pageSize int, eventStatus string, queryScope string) ([]*dto.EventListResponse, int, error)
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
	ListEventRegisteredUser(ctx context.Context, page, pageSize int, eventID int) ([]*dto.ListEventRegUserResponse, int, error)
}

// EventServiceImpl 实现 EventService 接口，提供事件相关的业务逻辑
type EventServiceImpl struct {
	eventRepo repository.EventRepository // 事件数据访问接口
	userRepo  userrepo.UserRepository    // 用户数据访问接口
	fileRepo  filerepo.FileRepository    // 文件数据访问接口
	msgSvc    msgsvc.MsgGroupService     // 消息群组服务接口
}

// NewEventService 创建服务实例
func NewEventService(
	eventRepo repository.EventRepository,
	userRepo userrepo.UserRepository,
	fileRepo filerepo.FileRepository,
	msgSvc msgsvc.MsgGroupService,
) EventService {
	return &EventServiceImpl{
		eventRepo: eventRepo,
		userRepo:  userRepo,
		fileRepo:  fileRepo,
		msgSvc:    msgSvc,
	}
}

// GetEventStatus 根据开始时间和结束时间计算活动状态
func (svc *EventServiceImpl) GetEventStatus(registrationStartTime time.Time, registrationEndTime time.Time) string {
	if registrationStartTime.After(time.Now()) {
		return "未开始"
	}
	if registrationStartTime.Before(time.Now()) && registrationEndTime.After(time.Now()) {
		return "正在进行"
	}
	if registrationEndTime.Before(time.Now()) {
		return "已结束"
	}
	return ""
}

// ListEvent 分页查询活动列表
func (svc *EventServiceImpl) ListEvent(ctx context.Context, page, pageSize int, eventStatus string, queryScope string) ([]*dto.EventListResponse, int, error) {
	return svc.eventRepo.List(ctx, page, pageSize, eventStatus, queryScope)
}

// GetEventDetail 获取活动详情
func (svc *EventServiceImpl) GetEventDetail(ctx context.Context, eventID int) (*model.Event, error) {
	event, err := svc.eventRepo.GetEventDetail(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// 获取关联图片列表
	images := svc.eventRepo.ListEventImage(ctx, eventID)

	// 添加图片到活动详情
	event.Images = make([]dto.Image, 0, len(images)) // 预分配空间，提高性能
	for _, img := range images {
		event.Images = append(event.Images, dto.Image{
			ImageID: img.ImageID,
			URL:     img.URL,
		})
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

	// 报名成功后，添加用户到活动对应的消息群组，消息群组添加失败不影响报名成功
	// 查询活动对应的消息群组
	group, count, err := svc.msgSvc.ListMsgGroups(ctx, 0, 0, "", eventID, "")
	if err != nil || count == 0 {
		// 不存在对应的消息群组，返回错误
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "进入活动消息群组失败，请联系管理员")
	}
	// 将用户添加到消息群组
	err = svc.msgSvc.AddUserToGroup(ctx, group[0].ID, []int{userID}, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "进入活动消息群组失败，请联系管理员")
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

	// 检查活动是否已开始
	if event.EventStartTime.Before(time.Now()) {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "活动已开始，无法取消报名")
	}

	// 执行取消报名逻辑
	err = svc.eventRepo.UpdateEUMapDeleteFlag(ctx, eventID, userID, utils.DeletedFlagYes)
	if err != nil {
		return err
	}

	// 取消报名成功后，将用户从活动对应的消息群组移除
	group, _, err := svc.msgSvc.ListMsgGroups(ctx, 0, 0, "", eventID, "")
	if err != nil {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "退出活动消息群组失败，请联系管理员处理")
	}
	// 将用户从消息群组移除
	err = svc.msgSvc.DeleteUserFromGroup(ctx, group[0].ID, []int{userID}, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "退出活动消息群组失败，请联系管理员处理")
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
	// 检查是否有重复的活动标题
	existingEvent, err := svc.eventRepo.GetEventByTitle(ctx, event.Title)
	if err != nil {
		return err
	}
	if existingEvent != nil {
		return utils.NewBusinessError(utils.ErrCodeResourceExists, "已存在同名活动，请修改标题后重试")
	}

	// 检查活动时间是否合理
	if event.EventStartTime.After(event.EventEndTime) {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "活动开始时间不能晚于结束时间")
	}
	if event.RegistrationStartTime.After(event.RegistrationEndTime) {
		return utils.NewBusinessError(utils.ErrCodeBusinessLogicError, "报名开始时间不能晚于结束时间")
	}

	// 使用 GORM 函数式事务
	err = svc.eventRepo.ExecTransaction(ctx, func(tx *gorm.DB) error {
		// 创建活动
		if err := svc.eventRepo.CreateEvent(ctx, tx, event); err != nil {
			return err
		}

		// 如果有图片，更新images表的biz_id和biz_type
		if len(imageIDList) > 0 {
			if err := svc.fileRepo.BatchUpdateImageBizID(ctx, tx, imageIDList, event.ID, utils.TypeEvent); err != nil {
				return err
			}
		}

		return nil // 返回 nil，GORM 自动提交
	})

	// 处理事务执行结果
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("事务执行失败: %w", err))
	}

	// 活动创建成功后，创建活动消息群组，消息群组创建失败不影响活动创建成功
	// 构建消息群组模型
	// 检查是否已存在对应的消息群组，理论上不应该存在
	_, count, _ := svc.msgSvc.ListMsgGroups(ctx, 0, 0, "", event.ID, "")
	if count > 0 {
		// 已存在对应的消息群组，直接返回成功，(在正确的业务流程下不应该出现这种情况)
		return nil
	}
	msgGroup := &msgmodel.UserMessageGroup{
		GroupName:      event.Title,
		Desc:           "由活动" + event.Title + "自动创建",
		EventID:        event.ID,
		IncludeAllUser: "N",
		CreateUser:     event.CreateUser,
		UpdateUser:     event.UpdateUser,
	}
	if err = svc.msgSvc.CreateMsgGroup(ctx, msgGroup, []int{}); err != nil {
		return utils.NewBusinessError(utils.ErrCodeServerInternalError, "自动创建活动消息群组失败"+err.Error())
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

	// 当标题修改时，检查是否有重复的活动标题
	if req.Title != nil && *req.Title != event.Title {
		existingEvent, err := svc.eventRepo.GetEventByTitle(ctx, *req.Title)
		if err != nil {
			return err
		}
		if existingEvent != nil {
			return utils.NewBusinessError(utils.ErrCodeResourceExists, "已存在同名活动，请修改标题后重试")
		}
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

	if len(updateFields) == 0 && len(imageIDList) == 0 {
		return nil // 无更新内容
	}

	// 设置更新人
	updateFields["update_user"] = userID

	// 使用 GORM 函数式事务
	err = svc.eventRepo.ExecTransaction(ctx, func(tx *gorm.DB) error {
		// 更新活动
		if err := svc.eventRepo.UpdateEvent(ctx, tx, eventID, updateFields); err != nil {
			return err
		}

		// 如果有图片，更新images表的biz_id和biz_type
		if len(imageIDList) > 0 {
			if err := svc.fileRepo.BatchUpdateImageBizID(ctx, tx, imageIDList, eventID, utils.TypeEvent); err != nil {
				return err
			}
		}

		return nil // 返回 nil，GORM 自动提交
	})

	// 处理事务执行结果
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("事务执行失败: %w", err))
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

	// 使用 GORM 函数式事务
	err = svc.eventRepo.ExecTransaction(ctx, func(tx *gorm.DB) error {
		// 软删除（更新is_deleted为Y，记录更新人）
		updateFields := map[string]interface{}{
			"is_deleted":  "Y",
			"update_user": userID,
		}
		if err := svc.eventRepo.UpdateEvent(ctx, tx, eventID, updateFields); err != nil {
			tx.Rollback()
			return err
		}

		return nil // 返回 nil，GORM 自动提交
	})

	// 处理事务执行结果
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("事务执行失败: %w", err))
	}
	// 删除活动成功后，删除活动对应的消息群组，消息群组删除失败不影响活动删除成功
	// 查询活动对应的消息群组
	group, count, err := svc.msgSvc.ListMsgGroups(ctx, 0, 0, "", eventID, "")
	if err != nil || count == 0 {
		// 不存在对应的消息群组，直接返回成功
		return nil
	}
	// 删除消息群组
	err = svc.msgSvc.DeleteMsgGroup(ctx, group[0].ID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ErrCodeServerInternalError, "删除活动消息群组失败"+err.Error())
	}
	return nil
}

// ListEventRegisteredUser 获取活动报名用户列表
func (svc *EventServiceImpl) ListEventRegisteredUser(ctx context.Context, page, pageSize int, eventID int) ([]*dto.ListEventRegUserResponse, int, error) {
	return svc.eventRepo.ListEventRegisteredUser(ctx, page, pageSize, eventID)
}
