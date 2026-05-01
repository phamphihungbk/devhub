package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDeploymentStatus = fmt.Errorf("invalid deployment status")
)

type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusRunning   DeploymentStatus = "running"
	DeploymentStatusCompleted DeploymentStatus = "completed"
	DeploymentStatusFailed    DeploymentStatus = "failed"
	// TODO: do we need this status, also update validation from usecase
	// DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
)

var deploymentStatusStringMapper = map[DeploymentStatus]string{
	DeploymentStatusPending:   "pending",
	DeploymentStatusRunning:   "running",
	DeploymentStatusCompleted: "completed",
	DeploymentStatusFailed:    "failed",
	// DeploymentStatusRolledBack: "rolled_back",
}

func (s DeploymentStatus) String() string {
	return deploymentStatusStringMapper[s]
}

func (s DeploymentStatus) IsValid() bool {
	switch s {
	case DeploymentStatusPending, DeploymentStatusRunning, DeploymentStatusCompleted, DeploymentStatusFailed:
		return true
	default:
		return false
	}
}

// Parse parses a string into a DeploymentStatus. It returns an error if the string is not a valid DeploymentStatus.
func (s DeploymentStatus) Parse(status string) (DeploymentStatus, error) {
	deploymentStatus := DeploymentStatus(status)

	if !deploymentStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidDeploymentStatus, status)
	}
	return deploymentStatus, nil
}

type Deployment struct {
	ID            uuid.UUID
	ServiceID     uuid.UUID
	EnvironmentID uuid.UUID
	Version       string
	Status        DeploymentStatus
	ExternalRef   string
	CommitSHA     string
	TriggeredBy   uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     time.Time
	FinishedAt    time.Time
}

type Deployments []Deployment
