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
	CreateImageFile(ctx context.Context, file *model.Image) error
}

// FileRepositoryImpl 文件存储库实现
type FileRepositoryImpl struct {
	db *gorm.DB
}

// NewFileRepository 创建文件存储库实例
func NewFileRepository(db *gorm.DB) FileRepository {
	return &FileRepositoryImpl{db: db}
}

// CreateImageFile 创建图片文件记录
func (repo *FileRepositoryImpl) CreateImageFile(ctx context.Context, image *model.Image) error {
	err := repo.db.WithContext(ctx).Create(image).Error
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("创建图片文件记录失败: %w", err))
	}
	return err
}
