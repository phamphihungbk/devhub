package config

import "time"

type ArgoCDConfig struct {
	Server            string
	AuthToken         string
	Insecure          bool
	Timeout           time.Duration
	ImageRegistryURL  string
	ImageRegistryHost string
	RepoBaseURL       string
	AutoCreateApp     bool
	AppProject        string
	AppNamespace      string
	AppDestServer     string
	AutoBuildImage    bool
	ImageBuilder      string
	MinikubeProfile   string
}

func LoadArgoCDConfig(cfg *ViperConfig) ArgoCDConfig {
	return ArgoCDConfig{
		Server:            cfg.GetString(ArgoCDServerKey),
		AuthToken:         cfg.GetString(ArgoCDAuthTokenKey),
		Insecure:          cfg.GetBool(ArgoCDInsecureKey),
		Timeout:           cfg.GetDuration(ArgoCDTimeoutKey),
		RepoBaseURL:       cfg.GetString(ArgoCDRepoBaseURLKey),
		ImageRegistryURL:  cfg.GetString(ArgoCDImageRegistryURLKey),
		ImageRegistryHost: cfg.GetString(ArgoCDImageRegistryHostKey),
		AutoCreateApp:     cfg.GetBool(ArgoCDAutoCreateAppKey),
		AppProject:        cfg.GetString(ArgoCDAppProjectKey),
		AppNamespace:      cfg.GetString(ArgoCDAppNamespaceKey),
		AppDestServer:     cfg.GetString(ArgoCDAppDestServerKey),
		AutoBuildImage:    cfg.GetBool(ArgoCDAutoBuildImageKey),
		ImageBuilder:      cfg.GetString(ArgoCDImageBuilderKey),
		MinikubeProfile:   cfg.GetString(ArgoCDMinikubeProfileKey),
	}
}
