package repository

import (
	"context"
	"fmt"
	"news-release/internal/file/model"
	"news-release/internal/utils"

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
	err := repo.db.WithContext(ctx).Create(file).Error
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("创建文件记录失败: %w", err))
	}
	return err
}
