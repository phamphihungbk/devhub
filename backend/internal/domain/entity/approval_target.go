package entity

import "fmt"

var (
	ErrInvalidApprovalResource = fmt.Errorf("invalid approval resource")
	ErrInvalidApprovalAction   = fmt.Errorf("invalid approval action")
)

type ApprovalResource string

const (
	ApprovalResourceScaffoldRequest ApprovalResource = "scaffold_request"
	ApprovalResourceDeployment      ApprovalResource = "deployment"
	ApprovalResourceRelease         ApprovalResource = "release"
)

func (r ApprovalResource) IsValid() bool {
	switch r {
	case ApprovalResourceScaffoldRequest, ApprovalResourceDeployment, ApprovalResourceRelease:
		return true
	default:
		return false
	}
}

func (r ApprovalResource) String() string {
	return string(r)
}

func (r ApprovalResource) Parse(resource string) (ApprovalResource, error) {
	parsed := ApprovalResource(resource)
	if !parsed.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidApprovalResource, resource)
	}
	return parsed, nil
}

type ApprovalAction string

const (
	ApprovalActionCreate ApprovalAction = "create"
	ApprovalActionUpdate ApprovalAction = "update"
	ApprovalActionDelete ApprovalAction = "delete"
)

func (a ApprovalAction) IsValid() bool {
	switch a {
	case ApprovalActionCreate, ApprovalActionUpdate, ApprovalActionDelete:
		return true
	default:
		return false
	}
}

func (a ApprovalAction) String() string {
	return string(a)
}

func (a ApprovalAction) Parse(action string) (ApprovalAction, error) {
	parsed := ApprovalAction(action)
	if !parsed.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidApprovalAction, action)
	}
	return parsed, nil
}

type ApprovalTarget struct {
	Resource ApprovalResource
	Action   ApprovalAction
}
