package controller

import (
	"net/http"
	"news-release/internal/notice/dto"
	"news-release/internal/notice/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
)

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
	// 初始化参数结构体并绑定查询参数
	var req dto.NoticeListRequest
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
	notice, total, err := n.noticeService.ListNotice(ctx, page, pageSize)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "服务器内部错误，调用服务层失败")
		return
	}

	var result []dto.NoticeResponse
	for _, n := range notice {
		result = append(result, dto.NoticeResponse{
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
