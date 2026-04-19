package config

// Token configuration environment variable keys
const (
	TokenDurationKey = "TOKEN_DURATION"
	TokenSecretKey   = "TOKEN_SECRET"
	TokenIssuerKey   = "TOKEN_ISSUER"
)

// Service configuration environment variable keys
const (
	ServiceNameKey      = "SERVICE_NAME"
	ServicePortKey      = "SERVICE_PORT"
	ServiceEnvKey       = "SERVICE_ENV"
	ServiceErrPrefixKey = "SERVICE_ERR_PREFIX"
	OtelExporterKey     = "OTEL_COLLECTOR_ENDPOINT"
	AdminAPIKey         = "ADMIN_API_KEY"    // #nosec G101
	AdminAPISecret      = "ADMIN_API_SECRET" // #nosec G101
	AppTimezoneKey      = "APP_TIMEZONE"
	SeatLockTTLKey      = "SEAT_LOCK_TTL"
)

// Database configuration environment variable keys
const (
	DatabaseURLKey             = "DATABASE_URL"
	DatabaseMaxOpenConnsKey    = "DATABASE_MAX_OPEN_CONNS"
	DatabaseMaxIdleConnsKey    = "DATABASE_MAX_IDLE_CONNS"
	DatabaseConnMaxLifetimeKey = "DATABASE_CONN_MAX_LIFETIME"  // duration string like "30m"
	DatabaseConnMaxIdleTimeKey = "DATABASE_CONN_MAX_IDLE_TIME" // duration string like "5m"
)

// Argo CD configuration environment variable keys
const (
	ArgoCDServerKey         = "ARGOCD_SERVER"
	ArgoCDAuthTokenKey      = "ARGOCD_AUTH_TOKEN"
	ArgoCDInsecureKey       = "ARGOCD_INSECURE"
	ArgoCDTimeoutKey        = "ARGOCD_TIMEOUT"
	ArgoCDAppNamespaceKey   = "ARGOCD_APP_NAMESPACE"
	ArgoCDTargetRevisionKey = "ARGOCD_TARGET_REVISION"
)

const (
	CIImageRegistryHostKey = "CI_IMAGE_REGISTRY_HOST"
	CIServerURLKey         = "CI_SERVER_URL"
)

// SCM configuration environment variable keys
const (
	SCMAPIURLKey      = "SCM_API_URL"
	SCMExternalURLKey = "SCM_EXTERNAL_URL"
	SCMInternalURLKey = "SCM_INTERNAL_URL"
	SCMTokenKey       = "SCM_TOKEN"
	SCMTimeoutKey     = "SCM_TIMEOUT"
)

// GitOps configuration environment variable keys
const (
	GitOpsRepoOwnerKey       = "GITOPS_REPO_OWNER"
	GitOpsRepoNameKey        = "GITOPS_REPO_NAME"
	GitOpsBranchKey          = "GITOPS_BRANCH"
	GitOpsBasePathKey        = "GITOPS_BASE_PATH"
	GitOpsCommitUserNameKey  = "GITOPS_COMMIT_USER_NAME"
	GitOpsCommitUserEmailKey = "GITOPS_COMMIT_USER_EMAIL"
	GitOpsTimeoutKey         = "GITOPS_TIMEOUT"
)

// Default configuration values
// These values are used if the environment variables are not set
var configDefaults = map[string]any{
	// Service configuration
	ServiceNameKey:      "devhub-backend-api",
	ServicePortKey:      ":8080",
	ServiceEnvKey:       "development",
	AppTimezoneKey:      "Asia/Hanoi",
	ServiceErrPrefixKey: "TR",
	SeatLockTTLKey:      "300s",

	// Database configuration
	DatabaseURLKey:             "postgres://devhub:devhubpass@devhub-db:5432/devhub?sslmode=disable",
	DatabaseMaxOpenConnsKey:    30,
	DatabaseMaxIdleConnsKey:    15,
	DatabaseConnMaxLifetimeKey: "30m",
	DatabaseConnMaxIdleTimeKey: "5m",

	// Argo CD configuration
	ArgoCDServerKey:         "argocd-server.argocd.svc.cluster.local:443",
	ArgoCDAuthTokenKey:      "your-token",
	ArgoCDInsecureKey:       true,
	ArgoCDTimeoutKey:        "10m",
	ArgoCDAppNamespaceKey:   "devhub",
	ArgoCDTargetRevisionKey: "main",

	// CI configuration
	CIImageRegistryHostKey: "host.docker.internal:5001",
	CIServerURLKey:         "http://host.docker.internal:3000",

	// SCM configuration
	SCMAPIURLKey:      "http://gitea:3000/api/v1",
	SCMExternalURLKey: "https://gitea.devhub.local",
	SCMInternalURLKey: "http://host.docker.internal:3000",
	SCMTimeoutKey:     "30s",

	// GitOps configuration
	GitOpsRepoOwnerKey:       "platform",
	GitOpsRepoNameKey:        "gitops-repo",
	GitOpsBranchKey:          "main",
	GitOpsBasePathKey:        "envs",
	GitOpsCommitUserNameKey:  "devhub-bot",
	GitOpsCommitUserEmailKey: "devhub-bot@local",
	GitOpsTimeoutKey:         "30s",

	// Token configuration
	TokenDurationKey: 3600,
	TokenSecretKey:   "your-secret-key",
	TokenIssuerKey:   "devhub-backend",
}
