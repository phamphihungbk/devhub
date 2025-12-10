package config

// Token configuration environment variable keys
const (
	TokenDurationKey = "TOKEN_DURATION"
	TokenSecretKey   = "TOKEN_SECRET"
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
	// Token configuration
	TokenDurationKey: 30,
	TokenSecretKey:   "your-secret-key",
}
