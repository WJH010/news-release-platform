package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局DB实例（初始化后复用，避免重复创建连接）
var db *gorm.DB

// NewDatabase 创建数据库连接
func NewDatabase(dsn string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 自动迁移模型（手动创建和修改数据库表结构时注释）
	// if err := migrateModels(db); err != nil {
	//  return nil, err
	// }

	return db, nil
}

// migrateModels 自动迁移数据库模型（GORM 的 AutoMigrate 方法会根据数据模型自动创建或更新表结构）
// func migrateModels(db *gorm.DB) error {
// 	// 添加需要迁移的模型
// 	return db.AutoMigrate(&model.Example{})
// }

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return db
}
