package notice

import (
	"context"
	"fmt"
	noticemodel "news-release/internal/model/notice"

	"gorm.io/gorm"
)

// 数据访问接口，定义数据访问的方法集
type NoticeRepository interface {
	// 分页查询政策列表
	List(ctx context.Context, page, pageSize int) ([]*noticemodel.Notice, int64, error)
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
func (r *NoticeRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*noticemodel.Notice, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var notices []*noticemodel.Notice
	query := r.db.WithContext(ctx)

	// 添加条件查询
	// 只展示有效公告
	query = query.Where("status = 1")

	// 按发布时间降序排列
	query = query.Order("release_time DESC")

	// 计算总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Model(&noticemodel.Notice{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("计算总数时数据库查询失败: %v", err)
	}

	// 查询数据
	if err := query.Offset(offset).Limit(pageSize).Find(&notices).Error; err != nil {
		return nil, 0, fmt.Errorf("数据库查询失败: %v", err)
	}

	return notices, total, nil
}
