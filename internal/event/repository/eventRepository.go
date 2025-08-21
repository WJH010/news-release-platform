package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"news-release/internal/event/model"
	usermodel "news-release/internal/user/model"
	"news-release/internal/utils"
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
	// IsUserRegistered 查询用户是否已报名活动
	IsUserRegistered(ctx context.Context, eventID int, userID int) (bool, error)
	// ListUserRegisteredEvents 获取用户已报名活动列表
	ListUserRegisteredEvents(ctx context.Context, page, pageSize int, userID int, eventStatus string) ([]*model.Event, int, error)
	// CreateEvent 创建活动
	CreateEvent(ctx context.Context, tx *gorm.DB, event *model.Event) error
	// UpdateEvent 更新活动
	UpdateEvent(ctx context.Context, tx *gorm.DB, eventID int, updateFields map[string]interface{}) error
	// DeleteEvent 删除活动
	DeleteEvent(ctx context.Context, eventID int, userID int) error
	// ListEventRegisteredUser 查询已报名活动的用户列表
	ListEventRegisteredUser(ctx context.Context, page, pageSize int, eventID int) ([]*usermodel.User, int, error)
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
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算总数时数据库查询失败: %v", err))
	}

	// 分页查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
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
			return nil, utils.NewBusinessError(utils.ErrCodeResourceNotFound, "活动不存在或已被删除，请刷新页面后重试")
		}
		return nil, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
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
	err := repo.db.WithContext(ctx).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		First(&mapping).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 映射不存在，只用来判断用户是否已报名活动，所以不返回异常
		}
		return nil, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return &mapping, nil
}

// CreatEventUserMap 创建活动-用户关联映射,将用户添加到活动中
func (repo *EventRepositoryImpl) CreatEventUserMap(ctx context.Context, eventUserMapping *model.EventUserMapping) error {
	err := repo.db.WithContext(ctx).Create(eventUserMapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return utils.NewBusinessError(utils.ErrCodeResourceExists, "已报名该活动，请勿重复报名")
		} else {
			return utils.NewSystemError(fmt.Errorf("创建活动-用户关联映射失败: %w", err))
		}
	}
	return nil
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
		return utils.NewSystemError(fmt.Errorf("数据更新异常: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "数据更新异常，未找到活动或状态已更新，请刷新页面后重试")
	}

	return nil
}

// IsUserRegistered 查询用户是否已报名活动
func (repo *EventRepositoryImpl) IsUserRegistered(ctx context.Context, eventID int, userID int) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&model.EventUserMapping{}).
		Where("event_id = ? AND user_id = ? AND is_deleted = ?", eventID, userID, "N").
		Count(&count).Error

	if err != nil {
		return false, utils.NewSystemError(fmt.Errorf("查询用户是否已报名活动失败: %w", err))
	}

	return count > 0, nil
}

// ListUserRegisteredEvents 获取用户已报名活动列表
func (repo *EventRepositoryImpl) ListUserRegisteredEvents(ctx context.Context, page, pageSize int, userID int, eventStatus string) ([]*model.Event, int, error) {
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

	query = query.Table("events e").
		Joins("JOIN event_user_mappings eum ON e.id = eum.event_id").
		Where("eum.user_id = ? AND e.is_deleted = ? AND eum.is_deleted = ?", userID, "N", "N").
		Find(&events)

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
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算总数时数据库查询失败: %v", err))
	}

	// 分页查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return events, int(total), nil
}

// CreateEvent 创建活动
func (repo *EventRepositoryImpl) CreateEvent(ctx context.Context, tx *gorm.DB, event *model.Event) error {
	// 插入新活动
	if err := tx.WithContext(ctx).Create(event).Error; err != nil {
		return utils.NewSystemError(fmt.Errorf("创建活动失败: %w", err))
	}

	return nil
}

// UpdateEvent 更新活动
func (repo *EventRepositoryImpl) UpdateEvent(ctx context.Context, tx *gorm.DB, eventID int, updateFields map[string]interface{}) error {
	// 更新活动信息
	result := tx.WithContext(ctx).Model(&model.Event{}).
		Where("id = ?", eventID).
		Updates(updateFields)

	if result.Error != nil {
		return utils.NewSystemError(fmt.Errorf("更新活动信息失败: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "更新活动信息失败，活动数据异常，请刷新页面后重试")
	}
	return nil
}

// DeleteEvent 删除活动
func (repo *EventRepositoryImpl) DeleteEvent(ctx context.Context, eventID int, userID int) error {
	// 软删除活动
	result := repo.db.WithContext(ctx).Model(&model.Event{}).
		Where("id = ?", eventID).
		Update("is_deleted", utils.DeletedFlagYes). // 采用软删除方式标记活动为已删除
		Update("update_user", userID)

	if result.Error != nil {
		return utils.NewSystemError(fmt.Errorf("软删除活动失败: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "删除活动失败，活动数据异常，请刷新页面后重试")
	}

	return nil
}

// ListEventRegisteredUser 查询已报名活动的用户列表
func (repo *EventRepositoryImpl) ListEventRegisteredUser(ctx context.Context, page, pageSize int, eventID int) ([]*usermodel.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var users []*usermodel.User
	var total int64

	query := repo.db.WithContext(ctx).
		Table("users u").
		Select("u.nickname, u.name, u.gender, u.phone_number, u.email, u.unit, u.department, u.position, u.industry, i.industry_name").
		Joins("JOIN event_user_mappings eum ON u.user_id = eum.user_id").
		Joins("LEFT JOIN industries i ON u.industry = i.industry_code").
		Where("eum.event_id = ? AND eum.is_deleted = ?", eventID, utils.DeletedFlagNo)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("计算总数时数据库查询失败: %v", err))
	}

	// 分页查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, utils.NewSystemError(fmt.Errorf("数据库查询失败: %v", err))
	}

	return users, int(total), nil
}
