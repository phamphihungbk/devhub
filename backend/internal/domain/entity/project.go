package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProjectEnvironment string

const (
	EnvDev  ProjectEnvironment = "dev"
	EnvProd ProjectEnvironment = "prod"
)

type Project struct {
	ID           uuid.UUID
	Name         string
	Description  string
	Environments []ProjectEnvironment
	CreatedBy    string
	DeletedAt    time.Time
}

type Projects []Project
