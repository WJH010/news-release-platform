package message

import (
	"errors"
	"fmt"
	"net/http"
	messagedto "news-release/internal/dto/message"
	msgsvc "news-release/internal/service/message"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 控制器
type MessageController struct {
	messageService msgsvc.MessageService
}

// 创建控制器实例
func NewMessageController(messageService msgsvc.MessageService) *MessageController {
	return &MessageController{messageService: messageService}
}

// 分页查询
func (m *MessageController) ListMessage(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req messagedto.MessageListRequest
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
	userID, exists := ctx.Get("userid")
	if !exists {
		utils.HandleError(ctx, nil, http.StatusInternalServerError, 0, "获取用户ID失败")
		return
	}
	// 类型转换
	uid, ok := userID.(int)
	if !ok {
		utils.HandleError(ctx, nil, http.StatusInternalServerError, 0, "用户ID类型错误")
		return
	}

	// 调用服务层
	message, total, err := m.messageService.ListMessage(ctx, page, pageSize, uid, req.MessageType)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []messagedto.MessageListResponse
	for _, m := range message {
		result = append(result, messagedto.MessageListResponse{
			ID:       m.ID,
			Title:    m.Title,
			Content:  m.Content,
			SendTime: m.SendTime,
			Type:     m.Type,
			TypeName: m.TypeName,
			Status:   m.Status,
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
	var req messagedto.MessageContentRequest
	if !utils.BindUrl(ctx, &req) {
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

	result := messagedto.MessageContentResponse{
		ID:       message.ID,
		Title:    message.Title,
		Content:  message.Content,
		SendTime: message.SendTime,
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
