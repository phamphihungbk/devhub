package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProjectEnvironment string

const (
	EnvDev     ProjectEnvironment = "dev"
	EnvProd    ProjectEnvironment = "prod"
	EnvStaging ProjectEnvironment = "staging"
)

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
