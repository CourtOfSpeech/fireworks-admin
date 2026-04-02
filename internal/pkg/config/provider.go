package config

import "github.com/google/wire"

// ProvideConfig 提供配置实例。
func ProvideConfig() (*Config, error) {
	return LoadByEnv()
}

// ProviderSet 配置依赖提供者集合。
var ProviderSet = wire.NewSet(ProvideConfig)
