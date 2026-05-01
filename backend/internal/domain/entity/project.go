package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidProjectEnvironment = fmt.Errorf("invalid project environment")
	ErrInvalidProjectStatus      = fmt.Errorf("invalid project status")
)

type Project struct {
	ID          uuid.UUID
	Name        string
	Description string
	OwnerTeamID uuid.UUID
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Projects []Project
