package config

import (
	"time"
)

type AppConfig struct {
	AdminAPIKey    string
	AdminAPISecret string
	Timezone       string
	SeatLockTTL    time.Duration
	// Add business feature flags here
}

func LoadAppConfig(cfg *ViperConfig) AppConfig {
	return AppConfig{
		AdminAPIKey:    cfg.GetString(AdminAPIKey),
		AdminAPISecret: cfg.GetString(AdminAPISecret),
		Timezone:       cfg.GetString(AppTimezoneKey),
		SeatLockTTL:    cfg.GetDuration(SeatLockTTLKey),
	}
}
