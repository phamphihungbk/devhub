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
	StatusPending DeploymentStatus = "pending"
	StatusRunning DeploymentStatus = "running"
	StatusSuccess DeploymentStatus = "success"
	StatusFailed  DeploymentStatus = "failed"
)

var deploymentStatusStringMapper = map[DeploymentStatus]string{
	StatusPending: "pending",
	StatusRunning: "running",
	StatusSuccess: "success",
	StatusFailed:  "failed",
}

func (s DeploymentStatus) String() string {
	return deploymentStatusStringMapper[s]
}

func (s DeploymentStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusRunning, StatusSuccess, StatusFailed:
		return true
	default:
		return false
	}
}

// Parse parses a string into a DeploymentStatus. It returns an error if the string is not a valid DeploymentStatus.
func (s DeploymentStatus) Parse(role string) (DeploymentStatus, error) {
	deploymentStatus := DeploymentStatus(role)

	if !deploymentStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidDeploymentStatus, role)
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
	TriggeredBy uuid.UUID
	CreatedAt   time.Time
}

type Deployments []Deployment
