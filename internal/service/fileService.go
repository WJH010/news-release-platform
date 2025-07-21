package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"news-release/internal/model"
	"news-release/internal/repository"
)

// FileService 文件服务接口
type FileService interface {
	UploadFile(ctx context.Context, fileHeader *FileHeader, articleID int, objectPrefix string, userID int) (*model.File, error)
}

// FileHeader 文件头信息
type FileHeader struct {
	OriginalFileName string
	ContentType      string
	Size             int64
	TemporaryFile    string
}

// FileServiceImpl 文件服务实现
type FileServiceImpl struct {
	minioRepo repository.MinIORepository
	fileRepo  repository.FileRepository
}

// NewFileService 创建文件服务实例
func NewFileService(minioRepo repository.MinIORepository, fileRepo repository.FileRepository) FileService {
	return &FileServiceImpl{
		minioRepo: minioRepo,
		fileRepo:  fileRepo,
	}
}

// UploadFile 上传文件
func (s *FileServiceImpl) UploadFile(ctx context.Context, fileHeader *FileHeader, articleID int, objectPrefix string, userID int) (*model.File, error) {
	// 生成唯一的对象名
	ext := filepath.Ext(fileHeader.OriginalFileName)
	objectName := fmt.Sprintf("%s/%d%s", objectPrefix, time.Now().UnixNano(), ext)

	// 上传到MinIO
	url, err := s.minioRepo.UploadFile(ctx, objectName, fileHeader.TemporaryFile)
	if err != nil {
		return nil, err
	}

	// 确定文件类型
	fileType := s.detectFileType(fileHeader.OriginalFileName)

	// 创建文件记录
	file := &model.File{
		ArticleID:    articleID,
		ObjectName:   objectName,
		URL:          url,
		FileName:     fileHeader.OriginalFileName,
		FileSize:     int(fileHeader.Size),
		ContentType:  fileHeader.ContentType,
		FileType:     string(fileType),
		UploadUserID: userID,
	}

	if err := s.fileRepo.CreateFile(ctx, file); err != nil {
		// 上传到数据库失败，删除MinIO中的文件
		_ = s.minioRepo.DeleteFile(ctx, objectName)
		return nil, err
	}

	return file, nil
}

// detectFileType 检测文件类型
func (s *FileServiceImpl) detectFileType(filename string) model.FileType {
	ext := strings.ToLower(filepath.Ext(filename))

	// 图片类型
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	for _, e := range imageExts {
		if ext == e {
			return model.FileTypeImage
		}
	}

	return model.FileTypeOther
}
