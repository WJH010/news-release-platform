package service

import (
	"context"
	"news-release/internal/notice/model"
	"news-release/internal/notice/repository"
)

// 服务接口，定义方法，接收 context.Context 和数据模型。
type NoticeService interface {
	ListNotice(ctx context.Context, page, pageSize int) ([]*model.Notice, int64, error)
}

// 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type NoticeServiceImpl struct {
	noticeRepo repository.NoticeRepository
}

// 创建服务实例
func NewNoticeService(noticeRepo repository.NoticeRepository) NoticeService {
	return &NoticeServiceImpl{noticeRepo: noticeRepo}
}

// 分页查询数据
func (s *NoticeServiceImpl) ListNotice(ctx context.Context, page, pageSize int) ([]*model.Notice, int64, error) {
	return s.noticeRepo.List(ctx, page, pageSize)
}
