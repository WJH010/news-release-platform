package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	Env   string
	Port  string
	DBURL string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	// .env文件
	// env := os.Getenv("ENV")
	env := "development" // 默认开发环境
	if env == "development" {
		if err := godotenv.Load("../.env"); err != nil {
			return nil, fmt.Errorf("加载开发环境配置失败: %v", err)
		}
	}

	return &Config{
		Env:   os.Getenv("ENV"),
		Port:  os.Getenv("PORT"),
		DBURL: os.Getenv("DB_URL"),
	}, nil
}
