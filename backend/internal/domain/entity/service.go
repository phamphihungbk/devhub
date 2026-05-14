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
	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Services []Service
