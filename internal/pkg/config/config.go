package config

import "strconv"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port            int      `mapstructure:"port"`
	Mode            string   `mapstructure:"mode"`
	AllowOrigins    []string `mapstructure:"allow_origins"`
	Timeout         int      `mapstructure:"timeout"`
	ReadTimeout     int      `mapstructure:"read_timeout"`
	WriteTimeout    int      `mapstructure:"write_timeout"`
	IdleTimeout     int      `mapstructure:"idle_timeout"`
	MaxHeaderBytes  int      `mapstructure:"max_header_bytes"`
	ShutdownTimeout int      `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	User             string `mapstructure:"user"`
	Password         string `mapstructure:"password"`
	DBName           string `mapstructure:"dbname"`
	SSLMode          string `mapstructure:"sslmode"`
	MaxOpenConns     int    `mapstructure:"max_open_conns"`
	MaxIdleConns     int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime  int    `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime  int    `mapstructure:"conn_max_idle_time"`
}

func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}

type CacheConfig struct {
	DefaultExpiration int `mapstructure:"default_expiration"`
	CleanupInterval   int `mapstructure:"cleanup_interval"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	Format    string `mapstructure:"format"`
	AddSource bool   `mapstructure:"add_source"`
}
