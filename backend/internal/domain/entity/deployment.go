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
	StatusPending    DeploymentStatus = "pending"
	StatusRunning    DeploymentStatus = "running"
	StatusCompleted  DeploymentStatus = "completed"
	StatusFailed     DeploymentStatus = "failed"
	StatusRolledBack DeploymentStatus = "rolled_back"
)

var deploymentStatusStringMapper = map[DeploymentStatus]string{
	StatusPending:    "pending",
	StatusRunning:    "running",
	StatusCompleted:  "completed",
	StatusFailed:     "failed",
	StatusRolledBack: "rolled_back",
}

func (s DeploymentStatus) String() string {
	return deploymentStatusStringMapper[s]
}

func (s DeploymentStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusRunning, StatusCompleted, StatusFailed, StatusRolledBack:
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
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Environment ProjectEnvironment
	Service     string
	Version     string
	Status      DeploymentStatus
	ExternalRef string
	CommitSHA   string
	TriggeredBy uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FinishedAt  *time.Time
}

type Deployments []Deployment
