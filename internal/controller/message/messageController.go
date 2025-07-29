package message

import (
	"errors"
	"net/http"
	msgsvc "news-release/internal/service/message"
	"news-release/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MessageListResponse 消息列表响应结构体
type MessageListResponse struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
	Type     string    `json:"type"`
	TypeName string    `json:"type_name"`
	Status   int       `json:"status"`
}

// MessageContentResponse 消息内容响应结构体
type MessageContentResponse struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
}

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
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	messageType := ctx.Query("message_type")
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

	// 转换参数类型
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用服务层
	message, total, err := m.messageService.ListMessage(ctx, page, pageSize, uid, messageType)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []MessageListResponse
	for _, m := range message {
		result = append(result, MessageListResponse{
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
	messageIDStr := ctx.Param("id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusBadRequest, 0, "消息ID无效")
		return
	}

	// 调用服务层
	message, err := m.messageService.GetMessageContent(ctx, messageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(ctx, err, http.StatusNotFound, 0, "消息不存在(id="+messageIDStr+")")
			return
		}
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "获取消息内容失败")
		return
	}

	result := MessageContentResponse{
		ID:       message.ID,
		Title:    message.Title,
		Content:  message.Content,
		SendTime: message.SendTime,
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
