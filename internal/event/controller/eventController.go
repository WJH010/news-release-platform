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
func (e *EventController) ListEvent(ctx *gin.Context) {
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
	event, total, err := e.eventService.ListEvent(ctx, page, pageSize, req.EventStatus)
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
			Images:                ev.Images,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      result,
	})
}
