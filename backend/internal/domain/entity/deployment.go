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
	DeploymentStatusPending    DeploymentStatus = "pending"
	DeploymentStatusRunning    DeploymentStatus = "running"
	DeploymentStatusCompleted  DeploymentStatus = "completed"
	DeploymentStatusFailed     DeploymentStatus = "failed"
	DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
)

var deploymentStatusStringMapper = map[DeploymentStatus]string{
	DeploymentStatusPending:    "pending",
	DeploymentStatusRunning:    "running",
	DeploymentStatusCompleted:  "completed",
	DeploymentStatusFailed:     "failed",
	DeploymentStatusRolledBack: "rolled_back",
}

func (s DeploymentStatus) String() string {
	return deploymentStatusStringMapper[s]
}

func (s DeploymentStatus) IsValid() bool {
	switch s {
	case DeploymentStatusPending, DeploymentStatusRunning, DeploymentStatusCompleted, DeploymentStatusFailed, DeploymentStatusRolledBack:
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
	ID           uuid.UUID
	ServiceID    uuid.UUID
	PluginID     uuid.UUID
	Environment  ProjectEnvironment
	Version      string
	Status       DeploymentStatus
	ExternalRef  string
	CommitSHA    string
	RunnerOutput string
	RunnerError  string
	TriggeredBy  uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	FinishedAt   *time.Time
}

type Deployments []Deployment
