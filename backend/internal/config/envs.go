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
	ArgoCDServerKey    = "ARGOCD_SERVER"
	ArgoCDAuthTokenKey = "ARGOCD_AUTH_TOKEN"
	ArgoCDInsecureKey  = "ARGOCD_INSECURE"
	ArgoCDTimeoutKey   = "ARGOCD_TIMEOUT"
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
	ArgoCDServerKey:    "argocd-server.argocd.svc.cluster.local:443",
	ArgoCDAuthTokenKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjphcGlLZXkiLCJuYmYiOjE3NzYwMjg1NTIsImlhdCI6MTc3NjAyODU1MiwianRpIjoiYWFlMTQ2NTQtYmRhZS00NzdlLWFlZjAtZGUyYTAxMzMyYjU2In0.2P5A4JDJY5543ajoy-PPJuzuiPQfS5rqKG7_ep4zBkk",
	ArgoCDInsecureKey:  true,
	ArgoCDTimeoutKey:   "10m",
	// Token configuration
	TokenDurationKey: 3600,
	TokenSecretKey:   "your-secret-key",
	TokenIssuerKey:   "devhub-backend",
}
