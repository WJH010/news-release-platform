package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	Env string
	// Port  string
	DBURL string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	var DBUrl string
	// 开发环境使用 .env 加载数据库配置
	env := os.Getenv("ENV")
	if env == "development" {
		if err := godotenv.Load("../.env"); err != nil {
			return nil, fmt.Errorf("加载开发环境配置失败: %v", err)
		}
		DBUrl = os.Getenv("DB_URL")
	} else {
		// 生产环境使用系统环境变量进行数据库配置注入
		DBUser := os.Getenv("DB_USER")
		DBPass := os.Getenv("DB_PASSWORD")
		DBHost := os.Getenv("DB_HOST")
		DBName := os.Getenv("DB_NAME")

		// 构建数据库连接字符串
		DBUrl = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
			DBUser, DBPass, DBHost, DBName)
	}

	return &Config{
		Env: os.Getenv("ENV"),
		// Port:  os.Getenv("PORT"),
		DBURL: DBUrl,
	}, nil
}
