package message

import (
	"net/http"
	msgsvc "news-release/internal/service/message"
	"news-release/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ArticleListResponse 文章列表响应结构体
type MessageListResponse struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
	Type     string    `json:"type"`
	Status   int       `json:"status"`
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
	message, total, err := m.messageService.ListMessage(ctx, page, pageSize, uid)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []MessageListResponse
	for _, p := range message {
		result = append(result, MessageListResponse{
			ID:       p.ID,
			Title:    p.Title,
			Content:  p.Content,
			SendTime: p.SendTime,
			Type:     p.Type,
			Status:   p.Status,
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
