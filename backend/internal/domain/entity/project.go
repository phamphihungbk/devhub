package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidProjectEnvironment = fmt.Errorf("invalid project environment")
)

type ProjectEnvironment string

const (
	EnvDev     ProjectEnvironment = "dev"
	EnvProd    ProjectEnvironment = "prod"
	EnvStaging ProjectEnvironment = "staging"
)

var projectEnvironmentStringMapper = map[ProjectEnvironment]string{
	EnvDev:     "dev",
	EnvProd:    "prod",
	EnvStaging: "staging",
}

func (s ProjectEnvironment) String() string {
	return projectEnvironmentStringMapper[s]
}

func (s ProjectEnvironment) IsValid() bool {
	switch s {
	case EnvDev, EnvProd, EnvStaging:
		return true
	default:
		return false
	}
}

// Parse parses a string into a ProjectEnvironment. It returns an error if the string is not a valid ProjectEnvironment.
func (s ProjectEnvironment) MustParse(env string) ProjectEnvironment {
	projectEnv := ProjectEnvironment(env)

	if !projectEnv.IsValid() {
		panic(`project environment: Parse(` + s + `): `)
	}
	return projectEnv
}

// Parse parses a string into a ProjectEnvironment. It returns an error if the string is not a valid ProjectEnvironment.
func (s ProjectEnvironment) Parse(env string) (ProjectEnvironment, error) {
	projectEnv := ProjectEnvironment(env)

	if !projectEnv.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidProjectEnvironment, env)
	}
	return projectEnv, nil
}

// ParseList parses a slice of strings into a slice of ProjectEnvironment.
// It returns an error if any string is not a valid ProjectEnvironment.
func (s ProjectEnvironment) ParseList(envs []string) ([]ProjectEnvironment, error) {
	var result []ProjectEnvironment
	for _, env := range envs {
		projectEnv := ProjectEnvironment(env)
		if !projectEnv.IsValid() {
			return nil, fmt.Errorf("%w: %s", ErrInvalidProjectEnvironment, env)
		}
		result = append(result, projectEnv)
	}
	return result, nil
}

type Project struct {
	ID           uuid.UUID
	Name         string
	Description  string
	Environments []ProjectEnvironment
	CreatedBy    uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type Projects []Project
