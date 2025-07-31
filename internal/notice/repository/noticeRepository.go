package repository

import (
	"context"
	"fmt"
	"news-release/internal/notice/model"

	"gorm.io/gorm"
)

// 数据访问接口，定义数据访问的方法集
type NoticeRepository interface {
	// 分页查询公告列表
	List(ctx context.Context, page, pageSize int) ([]*model.Notice, int64, error)
	// 根据id获取公告内容
	GetNoticeContent(ctx context.Context, noticeID int) (*model.Notice, error)
}

// 实现接口的具体结构体
type NoticeRepositoryImpl struct {
	db *gorm.DB
}

// 创建数据访问实例
func NewNoticeRepository(db *gorm.DB) NoticeRepository {
	return &NoticeRepositoryImpl{db: db}
}

// 分页查询数据
func (r *NoticeRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*model.Notice, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var notices []*model.Notice
	query := r.db.WithContext(ctx)

	// 添加条件查询
	// 只展示有效公告
	query = query.Where("status = 1")

	// 按发布时间降序排列
	query = query.Order("release_time DESC")

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Model(&model.Notice{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算总数时数据库查询失败: %v", err)
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&notices).Error; err != nil {
		return nil, 0, fmt.Errorf("数据库查询失败: %v", err)
	}

	return notices, total, nil
}

// 内容查询
func (r *NoticeRepositoryImpl) GetNoticeContent(ctx context.Context, noticeID int) (*model.Notice, error) {
	var notice model.Notice

	result := r.db.WithContext(ctx).First(&notice, noticeID)
	err := result.Error

	// 查询公告内容
	if err != nil {
		return nil, err
	}

	return &notice, nil
}
