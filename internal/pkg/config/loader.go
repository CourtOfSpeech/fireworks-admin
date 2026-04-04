package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadByEnv 根据环境变量 ENV 加载配置文件。
// 如果 ENV 未设置，默认使用 "dev"。
func LoadByEnv() (*Config, error) {
	env := "dev"
	if e := os.Getenv("ENV"); e != "" {
		env = e
	}

	configFile := fmt.Sprintf("configs/config.%s.toml", env)
	cfg, err := load(configFile)
	if err != nil {
		return nil, err
	}

	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// load 从指定路径加载配置文件。
// 支持环境变量覆盖，格式为 ${VAR:default}。
func load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("toml")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	resolveEnvVars(&cfg)

	return &cfg, nil
}

// resolveEnvVars 解析配置中的环境变量占位符。
// 优先使用环境变量值，未设置时使用默认值。
func resolveEnvVars(cfg *Config) {
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" && parsePort(v) > 0 {
		cfg.Database.Port = parsePort(v)
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		cfg.Database.SSLMode = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("APP_PORT"); v != "" && parsePort(v) > 0 {
		cfg.Server.Port = parsePort(v)
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
}

// parsePort 解析端口号字符串为整数。
func parsePort(s string) int {
	var port int
	fmt.Sscanf(s, "%d", &port)
	return port
}

// Validate 验证配置的必要字段。
// 生产环境下强制检查敏感配置项。
func Validate(cfg *Config) error {
	env := os.Getenv("ENV")
	isProduction := env == "production" || env == "prod"

	if isProduction {
		if cfg.Database.Password == "" || cfg.Database.Password == "postgres" {
			return errors.New("production environment requires a non-default database password via DB_PASSWORD")
		}
		if cfg.JWT.Secret == "" ||
			cfg.JWT.Secret == "your-secret-key-change-in-production" ||
			cfg.JWT.Secret == "default-secret-key-please-change-in-production" {
			return errors.New("production environment requires a secure JWT secret via JWT_SECRET")
		}
		if cfg.Database.SSLMode == "" || cfg.Database.SSLMode == "disable" {
			return errors.New("production environment requires SSL enabled for database connections")
		}
	}

	return nil
}
