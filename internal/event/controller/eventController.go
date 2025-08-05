package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"news-release/internal/event/dto"
	"news-release/internal/event/service"
	"news-release/internal/utils"
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
	event, total, err := ctr.eventService.ListEvent(ctx, page, pageSize, req.EventStatus)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "服务器内部错误，获取活动列表失败")
		return
	}

	var result []dto.EventListResponse
	for _, ev := range event {
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

// GetEventDetail 处理获取活动详情的请求
func (ctr *EventController) GetEventDetail(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.EventDetailRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 调用服务层获取活动详情
	event, err := ctr.eventService.GetEventDetail(ctx, req.EventID)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "服务器内部错误，获取活动详情失败")
		return
	}

	ctx.JSON(http.StatusOK, dto.EventDetailResponse{
		Title:                 event.Title,
		EventStartTime:        event.EventStartTime,
		EventEndTime:          event.EventEndTime,
		RegistrationStartTime: event.RegistrationStartTime,
		RegistrationEndTime:   event.RegistrationEndTime,
		EventAddress:          event.EventAddress,
		RegistrationFee:       event.RegistrationFee,
		Images:                event.Images,
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
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeAuthFailed, "获取用户ID失败")
		return
	}

	// 调用服务层进行活动报名
	err = ctr.eventService.RegistrationEvent(ctx, req.EventID, userID)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "服务器内部错误，活动报名失败")
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
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeAuthFailed, "获取用户ID失败")
		return
	}

	// 调用服务层查询用户是否报名该活动
	isRegistered, err := ctr.eventService.IsUserRegistered(ctx, req.EventID, userID)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, utils.ErrCodeServerInternalError, "服务器内部错误，查询用户报名状态失败")
		return
	}

	var flag, message string
	if isRegistered {
		flag = "Y"
		message = "已报名"
	} else {
		flag = "N"
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
