package entity

import "time"

type DeploymentStatus string

const (
	StatusPending DeploymentStatus = "pending"
	StatusRunning DeploymentStatus = "running"
	StatusSuccess DeploymentStatus = "success"
	StatusFailed  DeploymentStatus = "failed"
)

type Deployment struct {
	ID          string           `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID   string           `json:"projectId"`
	Environment string           `json:"environment"`
	Service     string           `json:"service"`
	Version     string           `json:"version"`
	Status      DeploymentStatus `gorm:"type:varchar(16)" json:"status"`
	TriggeredBy string           `json:"triggeredBy"`
	CreatedAt   time.Time        `json:"createdAt"`
}

type Deployments []Deployment
