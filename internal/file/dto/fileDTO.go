package dto

import "mime/multipart"

// FileUploadRequest 文件上传请求参数
type FileUploadRequest struct {
	BizType string                `form:"article_id" binding:"omitempty"` // 业务类型，ARTICLE-文章，EVENT-活动等
	BizID   int                   `form:"article_id" binding:"omitempty,numeric"`
	File    *multipart.FileHeader `form:"file" binding:"required"` // 文件字段，必填
}

type FileUploadResponse struct {
	ID       int    `json:"id"`
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}
