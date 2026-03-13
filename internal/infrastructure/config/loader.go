package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadByEnv loads the configuration file based on the environment variable ENV.
// If ENV is not set, it defaults to "dev".
func LoadByEnv() (*Config, error) {
	env := "dev"
	if e := viper.GetString("ENV"); e != "" {
		env = e
	}

	configFile := fmt.Sprintf("configs/config.%s.toml", env)
	return load(configFile)
}

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

	return &cfg, nil
}
