package entity

import "github.com/google/uuid"

type ScaffoldRequest struct {
	ID          uuid.UUID
	Template    string
	ProjectID   uuid.UUID
	Environment string
	Variables   map[string]string
}

type ScaffoldRequests []ScaffoldRequest
