package config

import "time"

type ArgoCDConfig struct {
	Server                 string
	AuthToken              string
	Insecure               bool
	Timeout                time.Duration
	AppProject             string
	AppNamespace           string
	TargetRevision         string
	RepoURL                string
	RepositoryRegistryHost string
}

func LoadArgoCDConfig(cfg *ViperConfig) ArgoCDConfig {
	return ArgoCDConfig{
		Server:                 cfg.GetString(ArgoCDServerKey),
		AuthToken:              cfg.GetString(ArgoCDAuthTokenKey),
		Insecure:               cfg.GetBool(ArgoCDInsecureKey),
		Timeout:                cfg.GetDuration(ArgoCDTimeoutKey),
		AppProject:             cfg.GetString(ArgoCDAppProjectKey),
		AppNamespace:           cfg.GetString(ArgoCDAppNamespaceKey),
		TargetRevision:         cfg.GetString(ArgoCDTargetRevisionKey),
		RepoURL:                cfg.GetString(ArgoCDTargetRevisionKey),
		RepositoryRegistryHost: cfg.GetString(ArgoCDRepositoryRegistryHost),
	}
}
