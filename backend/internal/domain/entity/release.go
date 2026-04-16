package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidReleaseStatus = fmt.Errorf("invalid release status")
)

type ReleaseStatus string

const (
	ReleaseStatusPending   ReleaseStatus = "pending"
	ReleaseStatusRunning   ReleaseStatus = "running"
	ReleaseStatusCompleted ReleaseStatus = "completed"
	ReleaseStatusFailed    ReleaseStatus = "failed"
)

var releaseStatusStringMapper = map[ReleaseStatus]string{
	ReleaseStatusPending:   "pending",
	ReleaseStatusRunning:   "running",
	ReleaseStatusCompleted: "completed",
	ReleaseStatusFailed:    "failed",
}

func (s ReleaseStatus) String() string {
	return releaseStatusStringMapper[s]
}

func (s ReleaseStatus) IsValid() bool {
	switch s {
	case ReleaseStatusPending, ReleaseStatusRunning, ReleaseStatusCompleted, ReleaseStatusFailed:
		return true
	default:
		return false
	}
}

// Parse parses a string into a DeploymentStatus. It returns an error if the string is not a valid DeploymentStatus.
func (s ReleaseStatus) Parse(status string) (ReleaseStatus, error) {
	releaseStatus := ReleaseStatus(status)

	if !releaseStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidReleaseStatus, status)
	}
	return releaseStatus, nil
}

type Release struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	PluginID    uuid.UUID
	Tag         string
	Target      string
	Name        string
	Notes       string
	HTMLURL     string
	ExternalRef string
	Status      ReleaseStatus
	TriggeredBy uuid.UUID
	CreatedAt   time.Time
}

type Releases []Release
