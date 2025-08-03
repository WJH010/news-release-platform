package model

import (
	"time"
)

// 活动状态常量定义
const (
	EventStatusInProgress = "InProgress" // 进行中
	EventStatusCompleted  = "Completed"  // 已结束
)

// Event 对应 events 表的数据模型
type Event struct {
	ID                    int       `json:"id" gorm:"primaryKey;column:id"`
	Title                 string    `json:"title" gorm:"type:varchar(255);not null;column:title"`               // 活动标题
	Detail                string    `json:"detail" gorm:"type:mediumtext;column:detail"`                        // 活动详情
	EventStartTime        time.Time `json:"event_start_time" gorm:"column:event_start_time"`                    // 活动开始时间
	EventEndTime          time.Time `json:"event_end_time" gorm:"column:event_end_time"`                        // 活动结束时间
	RegistrationStartTime time.Time `json:"registration_start_time" gorm:"column:registration_start_time"`      // 活动报名开始时间
	RegistrationEndTime   time.Time `json:"registration_end_time" gorm:"column:registration_end_time"`          // 活动报名截止时间
	EventAddress          string    `json:"event_address" gorm:"type:varchar(255);column:event_address"`        // 活动地址
	RegistrationFee       float64   `json:"registration_fee" gorm:"type:decimal(10,2);column:registration_fee"` // 报名费用
	CoverImageURL         string    `json:"cover_image_url" gorm:"column:cover_image_url"`                      // 封面图片URL
	IsDeleted             string    `json:"is_deleted" gorm:"column:is_deleted;default:N"`                      // 软删除标志
	CreateTime            time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`               // 数据创建时间，自动生成
	UpdateTime            time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`               // 数据最后更新时间，自动更新
	// 关联字段
	Images []string `json:"images" gorm:"-"` // 图片列表，存储图片URL
}

// TableName 设置表名
func (*Event) TableName() string {
	return "events"
}
