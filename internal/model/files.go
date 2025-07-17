package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// FileType 定义文件类型
type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeOther FileType = "other"
	// 后续可进行其他类型文件扩展，如video, document等
)

// File 文件模型
type File struct {
	ID          uint      `json:"id" gorm:"primaryKey;column:id"`
	ArticleType string    `json:"article_type" gorm:"column:article_type;type:varchar(255)"`
	ArticleID   int       `json:"article_id" gorm:"column:article_id;type:int"`
	ObjectName  string    `json:"object_name" gorm:"column:object_name;type:varchar(255)"`
	URL         string    `json:"url" gorm:"column:url;type:varchar(255)"`
	FileName    string    `json:"file_name" gorm:"column:file_name;type:varchar(255)"`
	FileSize    int       `json:"file_size" gorm:"column:file_size;type:int"`
	ContentType string    `json:"content_type" gorm:"column:content_type;type:varchar(255)"`
	FileType    string    `json:"file_type" gorm:"column:file_type;type:varchar(255)"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// Validate 验证数据
func (f *File) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}
