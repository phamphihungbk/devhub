package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type ViperOption func(v *viper.Viper) error

type ViperConfig struct {
	v *viper.Viper
}

type ConfigProvider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	All() map[string]interface{}
}

func NewViperConfig(opts ...ViperOption) (*ViperConfig, error) {
	v := viper.New()
	v.AutomaticEnv()

	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, err
		}
	}

	return &ViperConfig{v: v}, nil
}

func WithDefaults(defaults map[string]any) ViperOption {
	return func(v *viper.Viper) error {
		for key, val := range defaults {
			v.SetDefault(key, val)
		}
		return nil
	}
}

func WithOptionalConfigPaths(paths ...string) ViperOption {
	return func(v *viper.Viper) error {
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				v.SetConfigFile(path)
				if err := v.ReadInConfig(); err != nil {
					return fmt.Errorf("failed to read optional config file %s: %w", path, err)
				}
				break
			}
		}
		return nil
	}
}

func (c *ViperConfig) GetString(key string) string { return c.v.GetString(key) }

func (c *ViperConfig) GetInt(key string) int { return c.v.GetInt(key) }

func (c *ViperConfig) GetBool(key string) bool { return c.v.GetBool(key) }

func (c *ViperConfig) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *ViperConfig) All() map[string]interface{} {
	out := map[string]interface{}{}
	for _, key := range c.v.AllKeys() {
		out[key] = c.v.Get(key)
	}
	return out
}
