package config

import "time"

type ArgoCDConfig struct {
	Server      string
	AuthToken   string
	Insecure    bool
	Timeout     time.Duration
	RepoBaseURL string
}

func LoadArgoCDConfig(cfg *ViperConfig) ArgoCDConfig {
	return ArgoCDConfig{
		Server:      cfg.GetString(ArgoCDServerKey),
		AuthToken:   cfg.GetString(ArgoCDAuthTokenKey),
		Insecure:    cfg.GetBool(ArgoCDInsecureKey),
		Timeout:     cfg.GetDuration(ArgoCDTimeoutKey),
		RepoBaseURL: cfg.GetString(ArgoCDRepoBaseURLKey),
	}
}
