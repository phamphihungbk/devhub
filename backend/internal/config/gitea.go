package config

import "time"

type GiteaConfig struct {
	URL         string
	ExternalURL string
	Token       string
	Timeout     time.Duration
}

func LoadGiteaConfig(cfg *ViperConfig) GiteaConfig {
	return GiteaConfig{
		URL:         cfg.GetString(GiteaURLKey),
		ExternalURL: cfg.GetString(GiteaExternalURLKey),
		Token:       cfg.GetString(GiteaTokenKey),
		Timeout:     cfg.GetDuration(GiteaTimeoutKey),
	}
}
