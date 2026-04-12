package config

import (
	"log"
)

type Config struct {
	App     AppConfig      // Application-level settings such as API keys or feature flags
	ArgoCD  ArgoCDConfig   // Argo CD client settings used by deployment workers
	Service ServiceConfig  // Infrastructure-level service settings like name, port, and environment
	DB      DatabaseConfig // Database connection and pooling configuration
	Token   TokenConfig    // Token-related configuration
}

func MustConfigure() *Config {
	cfg, err := Configure()
	if err != nil {
		log.Fatalln(err)
	}
	return cfg
}

func Configure() (*Config, error) {
	cfg := MustConfig(
		WithOptionalConfigPaths("env.yaml", "../env.yaml"),
		WithDefaults(configDefaults),
	)

	return &Config{
		App:     LoadAppConfig(cfg),
		ArgoCD:  LoadArgoCDConfig(cfg),
		Service: LoadServiceConfig(cfg),
		DB:      LoadDatabaseConfig(cfg),
		Token:   LoadTokenConfig(cfg),
	}, nil
}

func MustConfig(opts ...ViperOption) *ViperConfig {
	cfg, err := NewViperConfig(opts...)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return cfg
}
