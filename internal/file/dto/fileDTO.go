package dto

import "mime/multipart"

// FileUploadRequest 文件上传请求参数
type FileUploadRequest struct {
	ArticleID int                   `form:"article_id" binding:"omitempty,numeric"`
	File      *multipart.FileHeader `form:"file" binding:"required"` // 文件字段，必填
}
