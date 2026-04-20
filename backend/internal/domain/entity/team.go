package entity

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID           uuid.UUID
	Name         string
	OwnerContact string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type Teams []Team
