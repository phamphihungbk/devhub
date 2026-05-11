package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEnvironmentTier   = fmt.Errorf("invalid environment tier")
	ErrInvalidEnvironmentConfig = fmt.Errorf("invalid environment config")
)

type EnvironmentConfig struct {
	Domain       string `json:"domain"`
	APIDomain    string `json:"api_domain"`
	Region       string `json:"region"`
	IngressClass string `json:"ingress_class"`
	PublicAccess bool   `json:"public_access"`
}

func (e EnvironmentConfig) Parse(variables string) (EnvironmentConfig, error) {

	var environmentConfig EnvironmentConfig
	err := json.Unmarshal([]byte(variables), &environmentConfig)

	if err != nil {
		return EnvironmentConfig{}, fmt.Errorf("%w: %s", ErrInvalidEnvironmentConfig, environmentConfig)
	}

	return environmentConfig, nil
}

func (e EnvironmentConfig) String() string {
	bytes, err := json.Marshal(e)

	if err != nil {
		return ""
	}

	return string(bytes)
}

type EnvironmentTier string

const (
	EnvironmentTierDevelopment EnvironmentTier = "development"
	EnvironmentTierStaging     EnvironmentTier = "staging"
	EnvironmentTierProduction  EnvironmentTier = "production"
)

func (t EnvironmentTier) IsValid() bool {
	switch t {
	case EnvironmentTierDevelopment, EnvironmentTierStaging, EnvironmentTierProduction:
		return true
	default:
		return false
	}
}

func (t EnvironmentTier) String() string {
	return string(t)
}

func (t EnvironmentTier) Parse(tier string) (EnvironmentTier, error) {
	parsed := EnvironmentTier(tier)
	if !parsed.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidEnvironmentTier, tier)
	}
	return parsed, nil
}

type Environment struct {
	ID             uuid.UUID
	ProjectID      uuid.UUID
	Name           string
	Tier           EnvironmentTier
	Cluster        string
	Namespace      string
	ArgocdInstance string
	Config         EnvironmentConfig
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Environments []Environment
