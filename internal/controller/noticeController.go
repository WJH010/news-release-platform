package controller

import (
	"net/http"
	"news-release/internal/service"
	"news-release/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NoticeResponse 公告列表响应结构体
type NoticeResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	ReleaseTime time.Time `json:"release_time"`
	Status      string    `json:"status"`
}

// 控制器
type NoticeController struct {
	noticeService service.NoticeService
}

// 创建控制器实例
func NewNoticeController(noticeService service.NoticeService) *NoticeController {
	return &NoticeController{noticeService: noticeService}
}

// 分页查询公告列表
func (n *NoticeController) ListNotice(ctx *gin.Context) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

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
	notice, total, err := n.noticeService.ListNotice(ctx, page, pageSize)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []NoticeResponse
	for _, n := range notice {
		result = append(result, NoticeResponse{
			ID:          n.ID,
			Title:       n.Title,
			Content:     n.Content,
			ReleaseTime: *n.ReleaseTime,
			Status:      map[int]string{1: "有效", 0: "无效"}[n.Status],
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
