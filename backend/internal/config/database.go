package config

import (
	"time"
)

// Database configuration keys
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func LoadDatabaseConfig(cfg *ViperConfig) DatabaseConfig {
	return DatabaseConfig{
		URL:             cfg.GetString(DatabaseURLKey),
		MaxOpenConns:    cfg.GetInt(DatabaseMaxOpenConnsKey),
		MaxIdleConns:    cfg.GetInt(DatabaseMaxIdleConnsKey),
		ConnMaxLifetime: cfg.GetDuration(DatabaseConnMaxLifetimeKey),
		ConnMaxIdleTime: cfg.GetDuration(DatabaseConnMaxIdleTimeKey),
	}
}
