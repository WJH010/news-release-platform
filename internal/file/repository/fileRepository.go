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
	// BatchUpdateImageBizID 批量更新图片的biz_id和biz_type
	BatchUpdateImageBizID(ctx context.Context, tx *gorm.DB, imageIDs []int, bizID int, bizType string) error
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

// BatchUpdateImageBizID 批量更新图片的biz_id和biz_type（新增）
func (repo *FileRepositoryImpl) BatchUpdateImageBizID(ctx context.Context, tx *gorm.DB, imageIDs []int, bizID int, bizType string) error {
	if len(imageIDs) == 0 {
		return nil
	}

	// 更新images表中指定ID的记录，设置biz_id和biz_type
	result := tx.WithContext(ctx).
		Table("images").
		Where("id IN (?)", imageIDs).
		Updates(map[string]interface{}{
			"biz_id":   bizID,
			"biz_type": bizType, // 标记为活动相关图片
		})

	if result.Error != nil {
		return utils.NewSystemError(fmt.Errorf("更新图片关联关系失败: %w", result.Error))
	}
	return nil
}
