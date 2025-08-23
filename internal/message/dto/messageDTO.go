package dto

import "time"

// MessageContentRequest 消息内容查询请求参数
type MessageContentRequest struct {
	MessageID int `uri:"id" binding:"required,numeric"` // 消息ID，必须为数字
}

// UnreadMessageCountRequest 获取未读消息数请求参数
type UnreadMessageCountRequest struct {
	MessageType string `form:"message_type"` // 消息类型
}

// ListUserGroupMessageRequest 消息群组列表请求参数
type ListUserGroupMessageRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`                       // 页码，最小为1
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`          // 页大小，1-100
	TypeCode string `form:"type_code" binding:"required,user_group_message_type"` // 消息类型代码
}

type GroupIDRequest struct {
	GroupID int `uri:"group_id" binding:"required,numeric"` // 消息组ID，必须为数字
}

// ListMessageByGroupsRequest 分页查询分组内消息列表请求参数
type ListMessageByGroupsRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`              // 页码，最小为1
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"` // 页大小，1-100
}

// MessageContentResponse 消息内容响应结构体
type MessageContentResponse struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
}

type MessageGroupDTO struct {
	MsgGroupID     int       `json:"msg_group_id"`
	GroupName      string    `json:"group_name"`
	LatestTitle    string    `json:"latest_title"`
	LatestContent  string    `json:"latest_content"`
	LatestSendTime time.Time `json:"latest_send_time"`
	HasUnread      string    `json:"has_unread"`
}

type ListMessageDTO struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
}
