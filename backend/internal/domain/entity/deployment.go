package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeploymentStatus string

const (
	StatusPending DeploymentStatus = "pending"
	StatusRunning DeploymentStatus = "running"
	StatusSuccess DeploymentStatus = "success"
	StatusFailed  DeploymentStatus = "failed"
)

type Deployment struct {
	ID          uuid.UUID
	ProjectID   string
	Environment string
	Service     string
	Version     string
	Status      DeploymentStatus
	TriggeredBy string
	CreatedAt   time.Time
}

type Deployments []Deployment
