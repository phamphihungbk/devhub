package entity

import (
	"time"

	"github.com/google/uuid"
)

type Release struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Tag         string
	Target      string
	Name        string
	Notes       string
	HTMLURL     string
	ExternalRef string
	TriggeredBy uuid.UUID
	CreatedAt   time.Time
}

type Releases []Release
