package config

import "time"

type ArgoCDConfig struct {
	Server          string
	AuthToken       string
	Insecure        bool
	Timeout         time.Duration
	RepoBaseURL     string
	AutoCreateApp   bool
	AppProject      string
	AppNamespace    string
	AppDestServer   string
	AutoBuildImage  bool
	ImageBuilder    string
	MinikubeProfile string
}

func LoadArgoCDConfig(cfg *ViperConfig) ArgoCDConfig {
	return ArgoCDConfig{
		Server:          cfg.GetString(ArgoCDServerKey),
		AuthToken:       cfg.GetString(ArgoCDAuthTokenKey),
		Insecure:        cfg.GetBool(ArgoCDInsecureKey),
		Timeout:         cfg.GetDuration(ArgoCDTimeoutKey),
		RepoBaseURL:     cfg.GetString(ArgoCDRepoBaseURLKey),
		AutoCreateApp:   cfg.GetBool(ArgoCDAutoCreateAppKey),
		AppProject:      cfg.GetString(ArgoCDAppProjectKey),
		AppNamespace:    cfg.GetString(ArgoCDAppNamespaceKey),
		AppDestServer:   cfg.GetString(ArgoCDAppDestServerKey),
		AutoBuildImage:  cfg.GetBool(ArgoCDAutoBuildImageKey),
		ImageBuilder:    cfg.GetString(ArgoCDImageBuilderKey),
		MinikubeProfile: cfg.GetString(ArgoCDMinikubeProfileKey),
	}
}
