package controller

import (
	"fmt"
	"net/http"
	"news-release/internal/file/dto"
	"news-release/internal/file/service"
	"news-release/internal/utils"
	"os"
	"path/filepath"

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
func (ctr *FileController) UploadFile(ctx *gin.Context) {
	// 初始化参数结构体并绑定查询参数
	var req dto.FileUploadRequest
	if !utils.BindForm(ctx, &req) {
		return
	}

	// 获取userID
	userID, err := utils.GetUserID(ctx)
	// 处理异常
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
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
	if err := ctx.SaveUploadedFile(req.File, tempFilePath); err != nil {
		err = utils.NewSystemError(fmt.Errorf("保存临时文件失败: %w", err))
		utils.WrapErrorHandler(ctx, err)
		return
	}
	defer os.Remove(tempFilePath)

	// 准备文件头信息
	fileHeader := &service.FileHeader{
		OriginalFileName: req.File.Filename,
		ContentType:      req.File.Header.Get("Content-Type"),
		Size:             req.File.Size,
		TemporaryFile:    tempFilePath,
	}

	// 上传文件
	response, err := ctr.fileService.UploadFile(ctx, fileHeader, req.BizType, req.BizID, userID)
	if err != nil {
		utils.WrapErrorHandler(ctx, err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
		"data":    response,
	})
}
