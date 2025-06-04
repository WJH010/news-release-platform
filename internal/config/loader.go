// config/loader.go 用于解析配置文件config.yaml
package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig 加载并解析配置文件
func LoadConfig(configPath string) (*Config, error) {

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 替换环境变量
	data = replaceEnvVariables(data)

	// 解析 YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// replaceEnvVariables 替换配置中的环境变量（格式: ${VAR_NAME} 或 ${VAR_NAME:-default}）
func replaceEnvVariables(data []byte) []byte {
	return []byte(os.Expand(string(data), func(key string) string {
		// 处理带默认值的情况: ${VAR:-default}
		if strings.Contains(key, ":-") {
			parts := strings.SplitN(key, ":-", 2)
			if val, ok := os.LookupEnv(parts[0]); ok {
				return val
			}
			return parts[1]
		}

		// 普通环境变量: ${VAR}
		return os.Getenv(key)
	}))
}

// validateConfig 验证配置的有效性
func validateConfig(config *Config) error {
	// 检查必要的数据库配置
	if config.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}
	if config.Database.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("数据库名不能为空")
	}

	return nil
}
