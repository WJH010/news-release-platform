package repository

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"news-release/internal/event/model"
	"time"
)

// EventRepository 数据访问接口，定义数据访问的方法集
type EventRepository interface {
	// List 分页查询
	List(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error)
	// ListEventImage 获取活动图片列表
	ListEventImage(ctx context.Context, bizIDs []int) []EventImage
}

// EventRepositoryImpl 实现接口的具体结构体
type EventRepositoryImpl struct {
	db *gorm.DB
}

// NewEventRepository 创建数据访问实例
func NewEventRepository(db *gorm.DB) EventRepository {
	return &EventRepositoryImpl{db: db}
}

// EventImage 结构体用于暂存图片查询结果,只在当前包内使用
type EventImage struct {
	BizID int    `json:"biz_id" gorm:"column:biz_id"`
	URL   string `json:"url" gorm:"column:url"`
}

// List 分页查询数据
func (r *EventRepositoryImpl) List(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var events []*model.Event
	var total int64

	query := r.db.WithContext(ctx)
	// 构建基础查询
	query = query.Table("events e").
		Where("is_deleted = ?", "N") // 软删除标志，查询未被删除的活动

	// 根据活动状态拼接查询条件
	if eventStatus == model.EventStatusInProgress {
		// 进行中的活动：报名时间在当前时间范围内
		query = query.Where("e.registration_start_time <= ? AND e.registration_end_time >= ?", time.Now(), time.Now())
		// 按活动开始时间升序排列
		query = query.Order("e.event_start_time ASC")
	} else if eventStatus == model.EventStatusCompleted {
		// 已结束的活动：报名截止时间在当前时间之前
		query = query.Where("e.registration_end_time < ?", time.Now())
		// 按活动开始时间降序排列
		query = query.Order("e.event_start_time DESC")
	}

	// 计算总数
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算总数时数据库查询失败: %w", err)
	}

	// 分页查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, int(total), nil
}

// ListEventImage 获取活动图片列表
func (r *EventRepositoryImpl) ListEventImage(ctx context.Context, bizIDs []int) []EventImage {
	var images []EventImage

	err := r.db.WithContext(ctx).
		Table("images").
		Where("biz_type = ? AND biz_id IN (?)", "EVENT", bizIDs).
		Find(&images).Error

	if err != nil {
		logrus.Errorf("获取活动图片失败: %v", err) // 只记录异常，不影响活动信息的返回
		return nil
	}

	return images
}
