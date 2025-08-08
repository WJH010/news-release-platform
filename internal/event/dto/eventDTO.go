package dto

import "time"

// EventListRequest 活动列表查询请求参数
type EventListRequest struct {
	Page        int    `form:"page" binding:"omitempty,min=1"`              // 页码，最小为1
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 页大小，1-100
	EventStatus string `form:"event_status"`                                // 活动状态
}

// EventDetailRequest 活动详情查询请求参数
type EventDetailRequest struct {
	EventID int `uri:"id" binding:"required,numeric"` // 活动ID，必须为数字
}

// EventRegistrationRequest 活动报名请求参数
type EventRegistrationRequest struct {
	EventID int `json:"event_id" binding:"required,numeric"` // 活动ID
}

// EventListResponse 活动列表响应结构体
type EventListResponse struct {
	ID                    int       `json:"id"`                      // 活动ID
	Title                 string    `json:"title"`                   // 活动标题
	EventStartTime        time.Time `json:"event_start_time"`        // 活动开始时间
	EventEndTime          time.Time `json:"event_end_time"`          // 活动结束时间
	RegistrationStartTime time.Time `json:"registration_start_time"` // 活动报名开始时间
	RegistrationEndTime   time.Time `json:"registration_end_time"`   // 活动报名截止时间
	EventAddress          string    `json:"event_address"`           // 活动地址
	RegistrationFee       float64   `json:"registration_fee"`        // 报名费用
	CoverImageURL         string    `json:"cover_image_url"`         // 封面图片URL
}

// EventDetailResponse 活动详情响应结构体
type EventDetailResponse struct {
	Title                 string    `json:"title"`                   // 活动标题
	Detail                string    `json:"detail"`                  // 活动内容
	EventStartTime        time.Time `json:"event_start_time"`        // 活动开始时间
	EventEndTime          time.Time `json:"event_end_time"`          // 活动结束时间
	RegistrationStartTime time.Time `json:"registration_start_time"` // 活动报名开始时间
	RegistrationEndTime   time.Time `json:"registration_end_time"`   // 活动报名截止时间
	EventAddress          string    `json:"event_address"`           // 活动地址
	RegistrationFee       float64   `json:"registration_fee"`        // 报名费用
	Images                []string  `json:"images"`                  // 图片列表
}
