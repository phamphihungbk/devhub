package config

type CIConfig struct {
	ImageRegistryHost string
	ServerURL         string
}

func LoadCIConfig(cfg *ViperConfig) CIConfig {
	return CIConfig{
		ImageRegistryHost: cfg.GetString(CIImageRegistryHostKey),
		ServerURL:         cfg.GetString(CIServerURLKey),
	}
}
