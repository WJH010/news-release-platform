package controller

import (
	"net/http"
	"news-release/internal/message/dto"
	"news-release/internal/message/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

// MessageController 控制器
type MessageController struct {
	messageService service.MessageService
}

// NewMessageController 创建控制器实例
func NewMessageController(messageService service.MessageService) *MessageController {
	return &MessageController{messageService: messageService}
}

// ListMessage 分页查询
func (ctr *MessageController) ListMessage(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.MessageListRequest
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

	// 调用服务层
	message, total, err := ctr.messageService.ListMessage(ctx, page, pageSize, userID, req.MessageType)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	var result []dto.MessageListResponse
	for _, m := range message {
		result = append(result, dto.MessageListResponse{
			ID:       m.ID,
			Title:    m.Title,
			Content:  m.Content,
			SendTime: m.SendTime,
			Type:     m.Type,
			TypeName: m.TypeName,
		})
	}

	// 返回分页结果
	ctx.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      result,
	})
}

// GetMessageContent 获取消息内容
func (ctr *MessageController) GetMessageContent(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.MessageContentRequest
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

	// 标记消息为已读
	err = ctr.messageService.MarkAsRead(ctx, userID, req.MessageID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层
	message, err := ctr.messageService.GetMessageContent(ctx, req.MessageID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	result := dto.MessageContentResponse{
		ID:       message.ID,
		Title:    message.Title,
		Content:  message.Content,
		SendTime: message.SendTime,
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GetUnreadMessageCount 获取未读消息数
func (ctr *MessageController) GetUnreadMessageCount(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.UnreadMessageCountRequest
	if !utils.BindQuery(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层
	count, err := ctr.messageService.GetUnreadMessageCount(ctx, userID, req.MessageType)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"count": count,
		},
	})
}

// MarkAllMessagesAsRead 一键已读
func (ctr *MessageController) MarkAllMessagesAsRead(ctx *gin.Context) {
	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 调用服务层
	err = ctr.messageService.MarkAllMessagesAsRead(ctx, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "所有消息已标记为已读",
	})
}
