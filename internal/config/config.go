package config

import (
	"time"
)

// Config 主配置结构
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	MinIO    MinIOConfig    `yaml:"minio"`
	Wechat   WechatConfig   `yaml:"wechat"` // 添加 Wechat 字段
	JWT      JWTConfig      `yaml:"jwt"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name           string   `yaml:"name"`
	Env            string   `yaml:"env"`
	Port           int      `yaml:"port"`
	Debug          bool     `yaml:"debug"`
	AllowedOrigins []string `yaml:"cors.allowed_origins"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver                string        `yaml:"driver"`
	Host                  string        `yaml:"host"`
	Port                  int           `yaml:"port"`
	Username              string        `yaml:"username"`
	Password              string        `yaml:"password"`
	DBName                string        `yaml:"dbname"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `yaml:"connection_max_lifetime"`
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	UseSSL          bool   `yaml:"use_ssl"`
	BucketName      string `yaml:"bucket_name"`
}

// WechatConfig 微信配置
type WechatConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	JwtSecret       string `yaml:"jwt_secret"`
	ExpirationHours int    `yaml:"expiration_hours"`
}
