package dto

import "mime/multipart"

type FileUploadRequest struct {
	ArticleID int                   `form:"article_id" binding:"omitempty,numeric"`
	File      *multipart.FileHeader `form:"file" binding:"required"` // 文件字段，必填
}
