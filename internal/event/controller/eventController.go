package controller

import (
	"net/http"
	"news-release/internal/event/dto"
	"news-release/internal/event/model"
	"news-release/internal/event/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

// EventController 定义事件控制器，处理与事件相关的 HTTP 请求
type EventController struct {
	eventService service.EventService // 事件服务接口
}

// NewEventController 创建事件控制器实例
func NewEventController(eventService service.EventService) *EventController {
	return &EventController{eventService: eventService}
}

// ListEvent 处理分页查询事件列表的请求
func (ctr *EventController) ListEvent(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.EventListRequest
	if !utils.BindQuery(ctx, &req) {
		return
	}

	// page 默认1
	page := req.Page
	if page == 0 {
		page = 1
	}

	// pageSize 默认10
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	// 调用服务层
	result, total, err := ctr.eventService.ListEvent(ctx, page, pageSize, req.EventStatus, req.QueryScope)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      result,
	})
}

// GetEventDetail 处理获取活动详情的请求
func (ctr *EventController) GetEventDetail(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.EventDetailRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 调用服务层获取活动详情
	event, err := ctr.eventService.GetEventDetail(ctx, req.EventID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 获取活动状态
	status := ctr.eventService.GetEventStatus(event.RegistrationStartTime, event.RegistrationEndTime)

	res := dto.EventDetailResponse{
		Title:                 event.Title,
		Detail:                event.Detail,
		EventStartTime:        event.EventStartTime,
		EventEndTime:          event.EventEndTime,
		RegistrationStartTime: event.RegistrationStartTime,
		RegistrationEndTime:   event.RegistrationEndTime,
		EventAddress:          event.EventAddress,
		RegistrationFee:       event.RegistrationFee,
		Status:                status,
		CoverImageURL:         event.CoverImageURL,
		Images:                event.Images,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

// RegistrationEvent 处理活动报名的请求
func (ctr *EventController) RegistrationEvent(ctx *gin.Context) {
	// 初始化参数结构体并绑定请求体
	var req dto.EventRegistrationRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层进行活动报名
	err = ctr.eventService.RegistrationEvent(ctx, req.EventID, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "活动报名成功",
	})
}

// IsUserRegistered 查询用户是否报名该活动
func (ctr *EventController) IsUserRegistered(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.EventDetailRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层查询用户是否报名该活动
	isRegistered, err := ctr.eventService.IsUserRegistered(ctx, req.EventID, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	var flag, message string
	if isRegistered {
		flag = utils.FlagYes
		message = "已报名"
	} else {
		flag = utils.FlagNo
		message = "未报名"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"is_registered": flag,
			"message":       message,
		},
	})
}

// CancelRegistrationEvent 处理取消活动报名的请求
func (ctr *EventController) CancelRegistrationEvent(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.EventDetailRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层取消活动报名
	err = ctr.eventService.CancelRegistrationEvent(ctx, req.EventID, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "取消活动报名成功",
	})
}

// ListUserRegisteredEvents 获取用户已报名的活动列表
func (ctr *EventController) ListUserRegisteredEvents(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.EventListRequest
	if !utils.BindQuery(ctx, &req) {
		return
	}

	// page 默认1
	page := req.Page
	if page == 0 {
		page = 1
	}

	// pageSize 默认10
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层获取用户已报名的活动列表
	events, total, err := ctr.eventService.ListUserRegisteredEvents(ctx, req.Page, req.PageSize, userID, req.EventStatus)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	var result []dto.EventListResponse
	for _, ev := range events {
		result = append(result, dto.EventListResponse{
			ID:                    ev.ID,
			Title:                 ev.Title,
			EventStartTime:        ev.EventStartTime,
			EventEndTime:          ev.EventEndTime,
			RegistrationStartTime: ev.RegistrationStartTime,
			RegistrationEndTime:   ev.RegistrationEndTime,
			EventAddress:          ev.EventAddress,
			RegistrationFee:       ev.RegistrationFee,
			CoverImageURL:         ev.CoverImageURL,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      result,
	})
}

// CreateEvent 处理创建活动的请求
func (ctr *EventController) CreateEvent(ctx *gin.Context) {
	// 初始化参数结构体并绑定请求体
	var req dto.CreateEventRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 转换时间格式
	eventStartTime, err := utils.StringToTime(req.EventStartTime)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	eventEndTime, err := utils.StringToTime(req.EventEndTime)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	registrationStartTime, err := utils.StringToTime(req.RegistrationStartTime)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	registrationEndTime, err := utils.StringToTime(req.RegistrationEndTime)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 构造活动模型
	event := &model.Event{
		Title:                 req.Title,
		Detail:                req.Detail,
		EventStartTime:        eventStartTime,
		EventEndTime:          eventEndTime,
		RegistrationStartTime: registrationStartTime,
		RegistrationEndTime:   registrationEndTime,
		EventAddress:          req.EventAddress,
		RegistrationFee:       req.RegistrationFee,
		CoverImageURL:         req.CoverImageURL,
		CreateUser:            userID,
		UpdateUser:            userID,
	}

	// 调用服务层创建活动
	err = ctr.eventService.CreateEvent(ctx, event, req.ImageIDList)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "活动创建成功",
		"data": gin.H{
			"event_id": event.ID,
		},
	})
}

// UpdateEvent 处理更新活动的请求
func (ctr *EventController) UpdateEvent(ctx *gin.Context) {
	// 获取活动ID
	var urlReq dto.EventDetailRequest
	if !utils.BindUrl(ctx, &urlReq) {
		return
	}
	// 初始化参数结构体并绑定请求体
	var req dto.UpdateEventRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层更新活动
	err = ctr.eventService.UpdateEvent(ctx, urlReq.EventID, req, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "活动更新成功",
	})
}

// DeleteEvent 处理删除活动的请求
func (ctr *EventController) DeleteEvent(ctx *gin.Context) {
	// 获取活动ID
	var req dto.EventDetailRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层删除活动
	err = ctr.eventService.DeleteEvent(ctx, req.EventID, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "活动删除成功",
	})
}

// ListEventRegisteredUsers 获取活动报名用户列表
func (ctr *EventController) ListEventRegisteredUsers(ctx *gin.Context) {
	// 获取活动ID
	var req dto.EventDetailRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// page 默认1
	page := 1

	// pageSize 默认10
	pageSize := 10

	// 调用服务层获取活动报名用户列表
	users, total, err := ctr.eventService.ListEventRegisteredUser(ctx, page, pageSize, req.EventID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      users,
	})
}
