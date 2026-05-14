package entity

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidJobType         = fmt.Errorf("invalid job type")
	ErrInvalidJobStatus       = fmt.Errorf("invalid job status")
	ErrInvalidJobResourceType = fmt.Errorf("invalid job resource type")
	ErrInvalidJobPayload      = fmt.Errorf("invalid job payload")
)

type JobPayload struct {
	Action        string                 `json:"action"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	ServiceID     string                 `json:"service_id,omitempty"`
	EnvironmentID string                 `json:"environment_id,omitempty"`
	Version       *string                `json:"version,omitempty"`
	TargetVersion *string                `json:"target_version,omitempty"`
}

func (j JobPayload) VariableString(key string) string {
	if j.Variables == nil {
		return ""
	}

	value, ok := j.Variables[key]
	if !ok || value == nil {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func (j JobPayload) Parse(variables string) (JobPayload, error) {

	var jobPayload JobPayload
	err := json.Unmarshal([]byte(variables), &jobPayload)

	if err != nil {
		return JobPayload{}, fmt.Errorf("%w: %s", ErrInvalidJobPayload, jobPayload)
	}

	return jobPayload, nil
}

func (j JobPayload) String() string {
	bytes, err := json.Marshal(j)

	if err != nil {
		return ""
	}

	return string(bytes)
}

type JobType string

const (
	JobTypeScaffold   JobType = "scaffold"
	JobTypeDeployment JobType = "deployment"
	JobTypeRelease    JobType = "release"
	JobTypeSyncEnv    JobType = "sync_env"
)

func (t JobType) IsValid() bool {
	switch t {
	case JobTypeScaffold, JobTypeDeployment, JobTypeRelease, JobTypeSyncEnv:
		return true
	default:
		return false
	}
}

func (t JobType) String() string {
	return string(t)
}

func (t JobType) Parse(jobType string) (JobType, error) {
	parsed := JobType(jobType)
	if !parsed.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidJobType, jobType)
	}
	return parsed, nil
}

type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusQueued, JobStatusRunning, JobStatusCompleted, JobStatusFailed:
		return true
	default:
		return false
	}
}

func (s JobStatus) String() string {
	return string(s)
}

func (s JobStatus) Parse(status string) (JobStatus, error) {
	parsed := JobStatus(status)
	if !parsed.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidJobStatus, status)
	}
	return parsed, nil
}

type JobResourceType string

const (
	JobResourceTypeScaffoldRequest JobResourceType = "scaffold_request"
	JobResourceTypeDeployment      JobResourceType = "deployment"
	JobResourceTypeRelease         JobResourceType = "release"
)

func (r JobResourceType) IsValid() bool {
	switch r {
	case JobResourceTypeScaffoldRequest, JobResourceTypeDeployment, JobResourceTypeRelease:
		return true
	default:
		return false
	}
}

func (r JobResourceType) String() string {
	return string(r)
}

func (r JobResourceType) Parse(resourceType string) (JobResourceType, error) {
	parsed := JobResourceType(resourceType)
	if !parsed.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidJobResourceType, resourceType)
	}
	return parsed, nil
}

type Job struct {
	ID           uuid.UUID
	Type         JobType
	Status       JobStatus
	ResourceType JobResourceType
	ResourceID   uuid.UUID
	PluginID     uuid.UUID
	Payload      JobPayload
	Result       string
	Error        string
	Attempts     int
	MaxAttempts  int
	CreatedBy    uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	StartedAt    time.Time
	FinishedAt   time.Time
}

type Jobs []Job
