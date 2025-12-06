package config

type ServiceConfig struct {
	Name         string
	Port         string
	Env          string
	ErrorPrefix  string
	OtelExporter string
}

func LoadServiceConfig(cfg *ViperConfig) ServiceConfig {
	return ServiceConfig{
		Name:         cfg.GetString(ServiceNameKey),
		Port:         cfg.GetString(ServicePortKey),
		Env:          cfg.GetString(ServiceEnvKey),
		ErrorPrefix:  cfg.GetString(ServiceErrPrefixKey),
		OtelExporter: cfg.GetString(OtelExporterKey),
	}
}
