package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"news-release/internal/file/model"
	"news-release/internal/file/repository"
)

// FileService 文件服务接口
type FileService interface {
	UploadFile(ctx context.Context, fileHeader *FileHeader, articleID int, userID int) (*model.File, error)
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
func (svc *FileServiceImpl) UploadFile(ctx context.Context, fileHeader *FileHeader, articleID int, userID int) (*model.File, error) {
	// 确定文件类型
	fileType := string(detectFileType(fileHeader.OriginalFileName))

	// 根据文件类型设置存储路径前缀
	objectPrefix := getObjectPrefixByType(fileType)

	// 生成唯一的对象名
	ext := filepath.Ext(fileHeader.OriginalFileName)
	objectName := fmt.Sprintf("%s/%d%s", objectPrefix, time.Now().UnixNano(), ext)

	// 上传到MinIO
	url, err := svc.minioRepo.UploadFile(ctx, objectName, fileHeader.TemporaryFile)
	if err != nil {
		return nil, err
	}

	// 创建文件记录
	file := &model.File{
		ArticleID:    articleID,
		ObjectName:   objectName,
		URL:          url,
		FileName:     fileHeader.OriginalFileName,
		FileSize:     int(fileHeader.Size),
		ContentType:  fileHeader.ContentType,
		FileType:     fileType,
		UploadUserID: userID,
	}

	if err := svc.fileRepo.CreateFile(ctx, file); err != nil {
		// 上传到数据库失败，删除MinIO中的文件
		_ = svc.minioRepo.DeleteFile(ctx, objectName)
		return nil, err
	}

	return file, nil
}

// detectFileType 检测文件类型
func detectFileType(filename string) model.FileType {
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

// getObjectPrefixByType 根据日期及文件类型获取存储路径前缀
func getObjectPrefixByType(fileType string) string {
	// 按月划分存储空间
	now := time.Now()
	yearMonth := now.Format("200601")

	switch fileType {
	case "image":
		return fmt.Sprintf("images/%s", yearMonth)
	default:
		return fmt.Sprintf("others/%s", yearMonth)
	}
}
