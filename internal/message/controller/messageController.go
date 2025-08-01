package controller

import (
	"errors"
	"fmt"
	"net/http"
	"news-release/internal/message/dto"
	"news-release/internal/message/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 控制器
type MessageController struct {
	messageService service.MessageService
}

// 创建控制器实例
func NewMessageController(messageService service.MessageService) *MessageController {
	return &MessageController{messageService: messageService}
}

// 分页查询
func (m *MessageController) ListMessage(ctx *gin.Context) {
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
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, err.Error())
		return
	}

	// 调用服务层
	message, total, err := m.messageService.ListMessage(ctx, page, pageSize, userID, req.MessageType)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
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

// 获取消息内容
func (m *MessageController) GetMessageContent(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.MessageContentRequest
	if !utils.BindUrl(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, err.Error())
		return
	}

	// 标记消息为已读
	err = m.messageService.MarkAsRead(ctx, userID, req.MessageID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "更新消息状态失败")
		return
	}

	// 调用服务层
	message, err := m.messageService.GetMessageContent(ctx, req.MessageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("文章不存在(id=%d)", req.MessageID)
			utils.HandleError(ctx, err, http.StatusNotFound, 0, msg)
			return
		}
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "获取消息内容失败")
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

// 获取未读消息数
func (m *MessageController) GetUnreadMessageCount(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.UnreadMessageCountRequest
	if !utils.BindQuery(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, err.Error())
		return
	}

	// 调用服务层
	count, err := m.messageService.GetUnreadMessageCount(ctx, userID, req.MessageType)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"count": count,
		},
	})
}
