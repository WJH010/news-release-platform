package dto

import "time"

// EventListRequest 活动列表查询请求参数
type EventListRequest struct {
	Page        int    `form:"page" binding:"omitempty,min=1"`              // 页码，最小为1
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 页大小，1-100
	EventStatus string `form:"event_status" binding:"omitempty"`            // 活动状态
	QueryScope  string `form:"query_scope" binding:"omitempty,query_scope"` // 查询范围，默认只查询未删除数据
}

// EventDetailRequest 活动详情查询请求参数
type EventDetailRequest struct {
	EventID int `uri:"id" binding:"required,numeric"` // 活动ID，必须为数字
}

// EventRegistrationRequest 活动报名请求参数
type EventRegistrationRequest struct {
	EventID int `json:"event_id" binding:"required,numeric"` // 活动ID
}

// CreateEventRequest 创建活动请求参数
type CreateEventRequest struct {
	Title                 string  `json:"title" binding:"required,max=255"`                       // 活动标题
	Detail                string  `json:"detail" binding:"required"`                              // 活动内容
	EventStartTime        string  `json:"event_start_time" binding:"required,time_format"`        // 活动开始时间
	EventEndTime          string  `json:"event_end_time" binding:"required,time_format"`          // 活动结束时间
	RegistrationStartTime string  `json:"registration_start_time" binding:"required,time_format"` // 活动报名开始时间
	RegistrationEndTime   string  `json:"registration_end_time" binding:"required,time_format"`   // 活动报名截止时间
	EventAddress          string  `json:"event_address" binding:"required,max=255"`               // 活动地址
	RegistrationFee       float64 `json:"registration_fee" binding:"gte=0"`                       // 报名费用，必须大于或等于 0
	CoverImageURL         string  `json:"cover_image_url" binding:"url"`                          // 封面图片URL
	ImageIDList           []int   `json:"image_id_list" binding:"omitempty,dive,min=1"`           // 图片ID列表
}

// UpdateEventRequest 更新活动请求参数
type UpdateEventRequest struct {
	Title                 *string  `json:"title" binding:"omitempty,non_empty_string,max=255"`                       // 活动标题
	Detail                *string  `json:"detail" binding:"omitempty,non_empty_string"`                              // 活动内容
	EventStartTime        *string  `json:"event_start_time" binding:"omitempty,non_empty_string,time_format"`        // 活动开始时间
	EventEndTime          *string  `json:"event_end_time" binding:"omitempty,non_empty_string,time_format"`          // 活动结束时间
	RegistrationStartTime *string  `json:"registration_start_time" binding:"omitempty,non_empty_string,time_format"` // 活动报名开始时间
	RegistrationEndTime   *string  `json:"registration_end_time" binding:"omitempty,non_empty_string,time_format"`   // 活动报名截止时间
	EventAddress          *string  `json:"event_address" binding:"omitempty,non_empty_string,max=255"`               // 活动地址
	RegistrationFee       *float64 `json:"registration_fee" binding:"omitempty,gte=0"`                               // 报名费用，必须大于或等于 0
	CoverImageURL         *string  `json:"cover_image_url" binding:"omitempty,url"`                                  // 封面图片URL
	ImageIDList           *[]int   `json:"image_id_list" binding:"omitempty,dive,min=1"`                             // 图片ID列表
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

// ListEventRegUserResponse 活动报名列表查询请求参数
type ListEventRegUserResponse struct {
	Nickname     string `json:"nickname"`
	Name         string `json:"name"`
	GenderCode   string `json:"gender_code"`
	Gender       string `json:"gender"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Unit         string `json:"unit"`
	Department   string `json:"department"`
	Position     string `json:"position"`
	Industry     int    `json:"industry"`
	IndustryName string `json:"industry_name"`
}
