// Package config 提供应用程序配置的结构定义和加载功能。
// 支持从 TOML 配置文件加载配置，并允许通过环境变量覆盖配置值。
// 配置包括服务器、数据库、JWT、缓存和日志等模块的设置。
package config

import "strconv"

// Config 应用程序主配置结构体。
// 包含服务器、数据库、JWT、缓存和日志等所有模块的配置。
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // 服务器配置
	Database DatabaseConfig `mapstructure:"database"` // 数据库配置
	JWT      JWTConfig      `mapstructure:"jwt"`      // JWT 认证配置
	Cache    CacheConfig    `mapstructure:"cache"`    // 缓存配置
	Log      LogConfig      `mapstructure:"log"`      // 日志配置
}

// ServerConfig HTTP 服务器配置结构体。
// 定义服务器的端口、运行模式、超时时间等参数。
type ServerConfig struct {
	Port            int      `mapstructure:"port"`             // 服务监听端口
	Mode            string   `mapstructure:"mode"`             // 运行模式：debug、release、test
	AllowOrigins    []string `mapstructure:"allow_origins"`    // CORS 允许的源列表
	Timeout         int      `mapstructure:"timeout"`          // 请求总超时时间（秒）
	ReadTimeout     int      `mapstructure:"read_timeout"`     // 读取请求超时时间（秒）
	WriteTimeout    int      `mapstructure:"write_timeout"`    // 写入响应超时时间（秒）
	IdleTimeout     int      `mapstructure:"idle_timeout"`     // 连接空闲超时时间（秒）
	MaxHeaderBytes  int      `mapstructure:"max_header_bytes"` // 请求头最大字节数
	StartTimeout    int      `mapstructure:"start_timeout"`    // 启动超时时间（秒）
	ShutdownTimeout int      `mapstructure:"shutdown_timeout"` // 优雅关闭超时时间（秒）
}

// DatabaseConfig 数据库连接配置结构体。
// 定义 PostgreSQL 数据库的连接参数和连接池设置。
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`              // 数据库主机地址
	Port            int    `mapstructure:"port"`              // 数据库端口
	User            string `mapstructure:"user"`              // 数据库用户名
	Password        string `mapstructure:"password"`          // 数据库密码
	DBName          string `mapstructure:"dbname"`            // 数据库名称
	SSLMode         string `mapstructure:"sslmode"`           // SSL 模式：disable、require、verify-ca、verify-full
	MaxOpenConns    int    `mapstructure:"max_open_conns"`    // 最大打开连接数
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 连接最大生命周期（秒）
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"` // 连接最大空闲时间（秒）
}

// DSN 生成 PostgreSQL 数据库连接字符串。
// 返回格式为：host=xxx port=xxx user=xxx password=xxx dbname=xxx sslmode=xxx
func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

// JWTConfig JWT 认证配置结构体。
// 定义 JWT 令牌的签名密钥和过期时间。
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`      // JWT 签名密钥
	ExpireTime int    `mapstructure:"expire_time"` // 令牌过期时间（秒）
}

// CacheConfig 缓存配置结构体。
// 定义内存缓存的默认过期时间和清理间隔。
type CacheConfig struct {
	DefaultExpiration int `mapstructure:"default_expiration"` // 默认过期时间（秒）
	CleanupInterval   int `mapstructure:"cleanup_interval"`   // 清理过期项的间隔（秒）
}

// LogConfig 日志配置结构体。
// 定义日志级别、输出格式和是否添加源码位置。
type LogConfig struct {
	Level     string `mapstructure:"level"`      // 日志级别：debug、info、warn、error
	Format    string `mapstructure:"format"`     // 输出格式：json、text
	AddSource bool   `mapstructure:"add_source"` // 是否在日志中添加源码位置
}
