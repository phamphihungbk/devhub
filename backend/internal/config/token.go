package config

type TokenConfig struct {
	Duration int
	Secret   string
}

func LoadTokenConfig(cfg *ViperConfig) TokenConfig {
	return TokenConfig{
		Duration: cfg.GetInt(TokenDurationKey),
		Secret:   cfg.GetString(TokenSecretKey),
	}
}
