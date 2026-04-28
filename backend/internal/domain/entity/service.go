package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	Name      string
	RepoURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Services []Service

type ServiceDependency struct {
	ID                 uuid.UUID
	ServiceID          uuid.UUID
	DependsOnServiceID uuid.UUID
	DependsOnService   *Service
	Type               string
	Protocol           string
	Port               *int
	Path               string
	Config             map[string]any
	CreatedBy          uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ServiceDependencies []ServiceDependency
