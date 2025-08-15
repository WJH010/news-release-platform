package controller

import (
	"github.com/gin-gonic/gin"
	"news-release/internal/message/dto"
	"news-release/internal/message/service"
	"news-release/internal/utils"
)

type MessageTypeController struct {
	messageTypeService service.MessageTypeService // 消息类型服务接口
}

// NewMessageTypeController 创建消息类型控制器实例
func NewMessageTypeController(messageTypeService service.MessageTypeService) *MessageTypeController {
	return &MessageTypeController{messageTypeService: messageTypeService}
}

// ListMessageType 列出所有消息类型
func (ctr *MessageTypeController) ListMessageType(ctx *gin.Context) {
	// 调用服务层获取消息类型列表
	messageTypes, err := ctr.messageTypeService.ListMessageType(ctx)
	if err != nil {
		// 处理异常
		utils.WrapErrorHandler(ctx, err)
		return
	}

	var list []dto.ListMessageTypeResponse
	for _, mt := range messageTypes {
		// 将消息类型的代码转换为字符串
		list = append(list, dto.ListMessageTypeResponse{
			ID:       mt.ID,
			TypeCode: mt.TypeCode,
			TypeName: mt.TypeName,
		})
	}

	// 返回结果
	ctx.JSON(200, gin.H{
		"data": list,
	})
}
