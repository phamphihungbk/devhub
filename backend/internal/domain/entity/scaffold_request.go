package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidScaffoldRequestVariables = fmt.Errorf("invalid scaffold request variables")
	ErrInvalidScaffoldRequestStatus    = fmt.Errorf("invalid scaffold request status")
)

type ScaffoldRequestVariables map[string]interface{}

// Parse parses a string into a ScaffoldRequestVariables. It returns an error if the string is not a valid ScaffoldRequestVariables.
func (s ScaffoldRequestVariables) Parse(variables string) (ScaffoldRequestVariables, error) {

	var scaffoldRequestVariables ScaffoldRequestVariables
	err := json.Unmarshal([]byte(variables), &scaffoldRequestVariables)

	if err != nil {
		return ScaffoldRequestVariables{}, fmt.Errorf("%w: %s", ErrInvalidScaffoldRequestVariables, variables)
	}

	return scaffoldRequestVariables, nil
}

func (e ScaffoldRequestVariables) String() string {
	bytes, err := json.Marshal(e)

	if err != nil {
		return ""
	}

	return string(bytes)
}

type ScaffoldRequestStatus string

const (
	ScaffoldRequestPending   ScaffoldRequestStatus = "pending"
	ScaffoldRequestApproved  ScaffoldRequestStatus = "approved"
	ScaffoldRequestRunning   ScaffoldRequestStatus = "running"
	ScaffoldRequestCompleted ScaffoldRequestStatus = "completed"
	ScaffoldRequestFailed    ScaffoldRequestStatus = "failed"
	// ScaffoldRequestRejected  ScaffoldRequestStatus = "rejected"
)

var scaffoldRequestStatusStringMapper = map[ScaffoldRequestStatus]string{
	ScaffoldRequestPending:   "pending",
	ScaffoldRequestApproved:  "approved",
	ScaffoldRequestRunning:   "running",
	ScaffoldRequestCompleted: "completed",
	ScaffoldRequestFailed:    "failed",
	// ScaffoldRequestRejected:  "rejected",
}

func (s ScaffoldRequestStatus) String() string {
	return scaffoldRequestStatusStringMapper[s]
}

func (s ScaffoldRequestStatus) IsValid() bool {
	switch s {
	case ScaffoldRequestPending, ScaffoldRequestApproved, ScaffoldRequestRunning, ScaffoldRequestCompleted, ScaffoldRequestFailed:
		return true
	default:
		return false
	}
}

// Parse parses a string into a ScaffoldRequestStatus. It returns an error if the string is not a valid ScaffoldRequestStatus.
func (s ScaffoldRequestStatus) Parse(status string) (ScaffoldRequestStatus, error) {
	scaffoldRequestStatus := ScaffoldRequestStatus(status)

	if !scaffoldRequestStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidScaffoldRequestStatus, scaffoldRequestStatus)
	}
	return scaffoldRequestStatus, nil
}

type ScaffoldRequest struct {
	ID            uuid.UUID
	ProjectID     uuid.UUID
	PluginID      uuid.UUID
	RequestedBy   uuid.UUID
	ApprovedBy    *uuid.UUID
	Status        ScaffoldRequestStatus
	Variables     ScaffoldRequestVariables
	ResultRepoURL string
	ApprovedAt    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ScaffoldRequests []ScaffoldRequest
