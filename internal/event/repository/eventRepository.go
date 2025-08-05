package repository

import (
	"context"
	"errors"
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
	// GetEventDetail 获取活动详情
	GetEventDetail(ctx context.Context, eventID int) (*model.Event, error)
	// ListEventImage 获取活动图片列表
	ListEventImage(ctx context.Context, bizID int) []EventImage
	// GetEventUserMap 查询活动-用户关联映射
	GetEventUserMap(ctx context.Context, eventID int, userID int) (*model.EventUserMapping, error)
	// CreatEventUserMap 创建活动-用户关联映射,将用户添加到活动中
	CreatEventUserMap(ctx context.Context, eventUserMapping *model.EventUserMapping) error
	// UpdateEUMapDeleteFlag 更新活动-用户关联删除标志
	UpdateEUMapDeleteFlag(ctx context.Context, eventID int, userID int, isDeleted string) error
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
func (repo *EventRepositoryImpl) List(ctx context.Context, page, pageSize int, eventStatus string) ([]*model.Event, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var events []*model.Event
	var total int64

	query := repo.db.WithContext(ctx)
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

// GetEventDetail 获取活动详情
func (repo *EventRepositoryImpl) GetEventDetail(ctx context.Context, eventID int) (*model.Event, error) {
	var event model.Event

	// 查询活动详情
	result := repo.db.WithContext(ctx).First(&event, eventID)
	err := result.Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("活动不存在或已被删除")
		}
		return nil, fmt.Errorf("获取活动详情失败: %w", err)
	}

	return &event, nil
}

// ListEventImage 获取活动图片列表
func (repo *EventRepositoryImpl) ListEventImage(ctx context.Context, bizID int) []EventImage {
	var images []EventImage

	err := repo.db.WithContext(ctx).
		Table("images").
		Where("biz_type = ? AND biz_id = ?", "EVENT", bizID).
		Find(&images).Error

	if err != nil {
		logrus.Errorf("获取活动图片失败: %v", err) // 只记录异常，不影响活动信息的返回
		return nil
	}

	return images
}

// GetEventUserMap 查询活动-用户关联映射
func (repo *EventRepositoryImpl) GetEventUserMap(ctx context.Context, eventID int, userID int) (*model.EventUserMapping, error) {
	var mapping model.EventUserMapping

	// 查询活动-用户关联映射
	result := repo.db.WithContext(ctx).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		First(&mapping)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 映射不存在
		}
		return nil, fmt.Errorf("查询活动-用户映射失败: %w", result.Error)
	}

	return &mapping, nil
}

// CreatEventUserMap 创建活动-用户关联映射,将用户添加到活动中
func (repo *EventRepositoryImpl) CreatEventUserMap(ctx context.Context, eventUserMapping *model.EventUserMapping) error {
	return repo.db.WithContext(ctx).Create(eventUserMapping).Error
}

// UpdateEUMapDeleteFlag 更新活动-用户关联删除标志
func (repo *EventRepositoryImpl) UpdateEUMapDeleteFlag(ctx context.Context, eventID int, userID int, isDeleted string) error {
	result := repo.db.WithContext(ctx).Model(&model.EventUserMapping{}).
		Where("event_id = ?", eventID).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_deleted": isDeleted,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
