package entity

import "github.com/google/uuid"

type ScaffoldVariables struct {
	ServiceName   string
	Port          int
	Database      string
	EnableLogging bool
}

type ScaffoldRequest struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Template    string
	Environment string
	Variables   ScaffoldVariables
}

type ScaffoldRequests []ScaffoldRequest
