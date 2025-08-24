package controller

import (
	"github.com/gin-gonic/gin"
	"news-release/internal/message/dto"
	"news-release/internal/message/model"
	"news-release/internal/message/service"
	"news-release/internal/utils"
)

type MsgGroupController struct {
	msgGroupService service.MsgGroupService
}

func NewMsgGroupController(msgGroupService service.MsgGroupService) *MsgGroupController {
	return &MsgGroupController{msgGroupService: msgGroupService}
}

// AddUserToGroup 用户入群
func (ctr *MsgGroupController) AddUserToGroup(ctx *gin.Context) {
	// 初始化参数结构体并绑定URL路径参数
	var urlReq dto.MsgGroupIDRequest
	if !utils.BindUrl(ctx, &urlReq) {
		return
	}
	// 初始化参数结构体并绑定查询参数
	var req dto.AddUserToGroupRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}
	// 获取当前登录userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	// 调用服务层
	err = ctr.msgGroupService.AddUserToGroup(ctx, urlReq.MsgGroupID, req.UserIDs, userID)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	// 返回成功响应
	ctx.JSON(200, gin.H{"message": "用户入群成功"})
}

// CreateMsgGroup 创建消息群组
func (ctr *MsgGroupController) CreateMsgGroup(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.CreateMsgGroupRequest
	if !utils.BindJSON(ctx, &req) {
		return
	}
	// 获取当前登录userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	// 构建消息群组模型
	msgGroup := &model.UserMessageGroup{
		GroupName:      req.GroupName,
		Desc:           req.Desc,
		IncludeAllUser: req.IncludeAllUser,
		CreateUser:     userID,
		UpdateUser:     userID,
	}
	// 调用服务层
	err = ctr.msgGroupService.CreateMsgGroup(ctx, msgGroup, req.UserIDs)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}
	// 返回成功响应
	ctx.JSON(200, gin.H{
		"message": "消息群组创建成功",
		"data": gin.H{
			"group_id": msgGroup.ID,
		},
	})
}
