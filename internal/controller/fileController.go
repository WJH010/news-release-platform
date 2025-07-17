package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"news-release/internal/service"
	"news-release/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileController 文件控制器
type FileController struct {
	fileService service.FileService
}

// NewFileController 创建文件控制器实例
func NewFileController(fileService service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// UploadFile 上传文件
func (c *FileController) UploadFile(ctx *gin.Context) {
	// 获取文章类型和ID
	articleType := ctx.PostForm("article_type")
	articleIDStr := ctx.PostForm("article_id")
	var articleID int
	// 转换 articleIDStr 参数
	if articleIDStr != "" {
		var err error
		articleID, err = strconv.Atoi(articleIDStr)

		if err != nil {
			utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "articleID格式转换错误")
			return
		}
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.HandleError(ctx, err, http.StatusBadRequest, 0, "未找到上传的文件")
		return
	}

	// 检查文件大小限制（示例：限制为10MB）
	// maxSize := int64(10 << 20) // 10MB
	// if file.Size > maxSize {
	// 	utils.HandleError(ctx, nil, http.StatusBadRequest, 0, fmt.Sprintf("文件大小超过限制（最大 %d MB）", maxSize/(1<<20)))
	// 	return
	// }

	// 保存临时文件
	tempFilePath := filepath.Join(os.TempDir(), uuid.New().String())
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "保存临时文件失败")
		return
	}
	defer os.Remove(tempFilePath)

	// 准备文件头信息
	fileHeader := &service.FileHeader{
		OriginalFileName: file.Filename,
		ContentType:      file.Header.Get("Content-Type"),
		Size:             file.Size,
		TemporaryFile:    tempFilePath,
	}

	// 根据文件类型设置存储路径前缀
	fileType := detectFileType(file.Filename, file.Header.Get("Content-Type"))
	objectPrefix := getObjectPrefixByType(fileType)

	// 上传文件
	fileInfo, err := c.fileService.UploadFile(ctx, fileHeader, articleType, articleID, objectPrefix)
	if err != nil {
		utils.HandleError(ctx, err, http.StatusInternalServerError, 0, "上传文件失败")
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "文件上传成功",
		"data":    fileInfo,
	})
}

// detectFileType 根据文件名和Content-Type检测文件类型
func detectFileType(filename, contentType string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// 图片类型
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	for _, e := range imageExts {
		if ext == e {
			return "image"
		}
	}

	return "other"
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
