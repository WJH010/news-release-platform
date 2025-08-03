package model

import (
	"time"
)

// Article 数据模型
type Article struct {
	ArticleID      int       `json:"article_id" gorm:"primaryKey;column:article_id"`
	ArticleTitle   string    `json:"article_title" gorm:"not null;column:article_title"`
	ArticleType    string    `json:"article_type" gorm:"not null;column:article_type"`
	ReleaseTime    time.Time `json:"release_time" gorm:"column:release_time"`
	Status         int       `json:"status" gorm:"column:status;default:1"` // 默认=1，1：有效
	BriefContent   string    `json:"brief_content" gorm:"type:text;column:brief_content"`
	ArticleContent string    `json:"article_content" gorm:"type:mediumtext;column:article_content"`
	IsSelection    int       `json:"is_selection" gorm:"default:2;column:is_selection"` // 默认=2，1：精选，2：非精选
	FieldID        int       `json:"field_id" gorm:"column:field_id"`
	CoverImageURL  string    `json:"cover_image_url" gorm:"column:cover_image_url"` // 封面图片URL
	ArticleSource  string    `json:"article_source" gorm:"column:article_source"`   // 文章来源
	CreateTime     time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime     time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	// 关联字段
	FieldName string `gorm:"-"` // 关联field_types
	TypeName  string `gorm:"-"` // 关联article_types

}

// TableName 设置表名
func (*Article) TableName() string {
	return "articles"
}
