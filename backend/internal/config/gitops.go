package config

import "time"

type SCMConfig struct {
	APIURL      string
	ExternalURL string
	Token       string
	Timeout     time.Duration
}

type GitOpsConfig struct {
	RepoOwner       string
	RepoName        string
	Branch          string
	BasePath        string
	CommitUserName  string
	CommitUserEmail string
	Timeout         time.Duration
}

func LoadSCMConfig(cfg *ViperConfig) SCMConfig {
	return SCMConfig{
		APIURL:      cfg.GetString(SCMAPIURLKey),
		ExternalURL: cfg.GetString(SCMExternalURLKey),
		Token:       cfg.GetString(SCMTokenKey),
		Timeout:     cfg.GetDuration(SCMTimeoutKey),
	}
}

func LoadGitOpsConfig(cfg *ViperConfig) GitOpsConfig {
	return GitOpsConfig{
		RepoOwner:       cfg.GetString(GitOpsRepoOwnerKey),
		RepoName:        cfg.GetString(GitOpsRepoNameKey),
		Branch:          cfg.GetString(GitOpsBranchKey),
		BasePath:        cfg.GetString(GitOpsBasePathKey),
		CommitUserName:  cfg.GetString(GitOpsCommitUserNameKey),
		CommitUserEmail: cfg.GetString(GitOpsCommitUserEmailKey),
		Timeout:         cfg.GetDuration(GitOpsTimeoutKey),
	}
}
