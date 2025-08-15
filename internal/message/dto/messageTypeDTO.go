package dto

// ListMessageTypeResponse 消息类型列表响应结构体
type ListMessageTypeResponse struct {
	ID       int    `json:"id"`
	TypeCode string `json:"type_code"`
	TypeName string `json:"type_name"`
}
