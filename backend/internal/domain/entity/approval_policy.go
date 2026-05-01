package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApprovalPolicy struct {
	ID                uuid.UUID
	Resource          string
	Action            string
	ProjectID         *uuid.UUID
	ServiceID         *uuid.UUID
	EnvironmentID     uuid.UUID
	RequiredApprovals int
	Enabled           bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ApprovalPolicies []ApprovalPolicy
