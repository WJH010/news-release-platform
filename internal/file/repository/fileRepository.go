package repository

import (
	"context"
	"news-release/internal/file/model"

	"gorm.io/gorm"
)

// FileRepository 文件存储库接口
type FileRepository interface {
	CreateFile(ctx context.Context, file *model.File) error
}

// FileRepositoryImpl 文件存储库实现
type FileRepositoryImpl struct {
	db *gorm.DB
}

// NewFileRepository 创建文件存储库实例
func NewFileRepository(db *gorm.DB) FileRepository {
	return &FileRepositoryImpl{db: db}
}

// CreateFile 创建文件记录
func (repo *FileRepositoryImpl) CreateFile(ctx context.Context, file *model.File) error {
	return repo.db.WithContext(ctx).Create(file).Error
}
