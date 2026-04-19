package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ScaffoldRequestVariables struct {
	ServiceName   string `json:"service_name"`
	ModulePath    string `json:"module_path"`
	Port          int    `json:"port"`
	Database      string `json:"database"`
	EnableLogging bool   `json:"enable_logging"`
}

// Parse parses a string into a ProjectEnvironment. It returns an error if the string is not a valid ProjectEnvironment.
func (s ScaffoldRequestVariables) Parse(variables string) (ScaffoldRequestVariables, error) {

	var scaffoldRequestVariables ScaffoldRequestVariables
	err := json.Unmarshal([]byte(variables), &scaffoldRequestVariables)

	if err != nil {
		return ScaffoldRequestVariables{}, fmt.Errorf("%w: %s", ErrInvalidProjectEnvironment, scaffoldRequestVariables)
	}

	return scaffoldRequestVariables, nil
}

func (s ScaffoldRequestVariables) String() string {
	bytes, err := json.Marshal(s)

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
	ScaffoldRequestRejected  ScaffoldRequestStatus = "rejected"
)

var scaffoldRequestStatusStringMapper = map[ScaffoldRequestStatus]string{
	ScaffoldRequestPending:   "pending",
	ScaffoldRequestApproved:  "approved",
	ScaffoldRequestRunning:   "running",
	ScaffoldRequestCompleted: "completed",
	ScaffoldRequestFailed:    "failed",
	ScaffoldRequestRejected:  "rejected",
}

func (s ScaffoldRequestStatus) String() string {
	return scaffoldRequestStatusStringMapper[s]
}

func (s ScaffoldRequestStatus) IsValid() bool {
	switch s {
	case ScaffoldRequestPending, ScaffoldRequestApproved, ScaffoldRequestRunning, ScaffoldRequestCompleted, ScaffoldRequestFailed, ScaffoldRequestRejected:
		return true
	default:
		return false
	}
}

// Parse parses a string into a ScaffoldRequestStatus. It returns an error if the string is not a valid ScaffoldRequestStatus.
func (s ScaffoldRequestStatus) Parse(status string) (ScaffoldRequestStatus, error) {
	scaffoldRequestStatus := ScaffoldRequestStatus(status)

	if !scaffoldRequestStatus.IsValid() {
		return "", fmt.Errorf("invalid scaffold request status: %s", status)
	}
	return scaffoldRequestStatus, nil
}

type ScaffoldRequest struct {
	ID            uuid.UUID
	PluginID      uuid.UUID
	ProjectID     uuid.UUID
	RequestedBy   uuid.UUID
	Status        ScaffoldRequestStatus
	Environment   ProjectEnvironment
	Variables     ScaffoldRequestVariables
	ApprovedBy    *uuid.UUID
	ResultRepoURL string
	ApprovedAt    *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ScaffoldRequests []ScaffoldRequest
