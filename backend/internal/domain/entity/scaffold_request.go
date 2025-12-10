package entity

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ScaffoldRequestVariables struct {
	ServiceName   string
	Port          int
	Database      string
	EnableLogging bool
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

type ScaffoldRequest struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Template    string
	Environment ProjectEnvironment
	Variables   ScaffoldRequestVariables
}

type ScaffoldRequests []ScaffoldRequest
