package dto

import "time"

// ListArticleRequest 文章列表查询请求参数
type MessageListRequest struct {
	Page        int    `form:"page" binding:"omitempty,min=1"`              // 页码，最小为1
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 页大小，1-100
	MessageType string `form:"message_type"`                                // 消息类型
}

// 消息内容查询请求参数
type MessageContentRequest struct {
	MessageID int `form:"id" binding:"required,numeric"` // 消息ID，必须为数字
}

// 获取未读消息数请求参数
type UnreadMessageCountRequest struct {
	MessageType string `form:"message_type"` // 消息类型
}

// MessageListResponse 消息列表响应结构体
type MessageListResponse struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
	Type     string    `json:"type"`
	TypeName string    `json:"type_name"`
}

// MessageContentResponse 消息内容响应结构体
type MessageContentResponse struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
}
