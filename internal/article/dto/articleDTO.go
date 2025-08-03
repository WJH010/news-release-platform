package dto

import "time"

// ArticleListRequest 文章列表查询请求参数
type ArticleListRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`               // 页码，最小为1
	PageSize     int    `form:"page_size" binding:"omitempty,min=1,max=100"`  // 页大小，1-100
	ArticleTitle string `form:"article_title"`                                // 文章标题
	FieldID      int    `form:"field_id" binding:"omitempty,numeric"`         // 领域ID，必须为数字
	IsSelection  int    `form:"is_selection" binding:"omitempty,numeric"`     // 是否精选，必须为数字
	ArticleType  string `form:"article_type"`                                 // 文章类型
	ReleaseTime  string `form:"release_time" binding:"omitempty,time_format"` // 发布时间
	Status       int    `form:"status" binding:"omitempty,numeric"`           // 状态，必须为数字
}

// ArticleContentRequest 文章内容查询请求参数
type ArticleContentRequest struct {
	ArticleID int `uri:"id" binding:"required,numeric"` // 文章ID，必须为数字
}

// ArticleListResponse 文章列表响应结构体
type ArticleListResponse struct {
	ArticleID       int       `json:"article_id"`
	ArticleTitle    string    `json:"article_title"`
	ArticleTypeCode string    `json:"article_type_code"`
	ArticleType     string    `json:"article_type"`
	FieldName       string    `json:"field_name"`
	ReleaseTime     time.Time `json:"release_time"`
	BriefContent    string    `json:"brief_content"`
	IsSelection     int       `json:"is_selection"`
	CoverImageURL   string    `json:"cover_image_url"`
	ArticleSource   string    `json:"article_source"`
}

// ArticleContentResponse 文章内容响应结构体
type ArticleContentResponse struct {
	ArticleID       int       `json:"article_id"`
	ArticleTitle    string    `json:"article_title"`
	FieldName       string    `json:"field_name"`
	ReleaseTime     time.Time `json:"release_time"`
	ArticleContent  string    `json:"article_content"`
	ArticleTypeCode string    `json:"article_type_code"`
	ArticleType     string    `json:"article_type"`
	ArticleSource   string    `json:"article_source"`
}
