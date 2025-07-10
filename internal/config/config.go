package config

import (
	"time"
)

// Config 主配置结构
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
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

// WechatConfig 微信配置
type WechatConfig struct {
	AppID  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	SecretKey       string `yaml:"secret_key"`
	ExpirationHours int    `yaml:"expiration_hours"`
}
