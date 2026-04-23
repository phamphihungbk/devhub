package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidApprovalRequestStatus = fmt.Errorf("invalid approval request status")
	ErrInvalidApprovalDecision      = fmt.Errorf("invalid approval decision")
	ErrInvalidApprovalResource      = fmt.Errorf("invalid approval resource")
	ErrInvalidApprovalAction        = fmt.Errorf("invalid approval action")
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

func ParseOptionalApprovalResource(resource string) (*ApprovalResource, error) {
	if resource == "" {
		return nil, nil
	}

	parsed, err := new(ApprovalResource).Parse(resource)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
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

func ParseOptionalApprovalAction(action string) (*ApprovalAction, error) {
	if action == "" {
		return nil, nil
	}

	parsed, err := new(ApprovalAction).Parse(action)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

type ApprovalTarget struct {
	Resource ApprovalResource
	Action   ApprovalAction
}

type ApprovalRequestStatus string

const (
	ApprovalRequestStatusPending  ApprovalRequestStatus = "pending"
	ApprovalRequestStatusApproved ApprovalRequestStatus = "approved"
	ApprovalRequestStatusRejected ApprovalRequestStatus = "rejected"
	ApprovalRequestStatusCanceled ApprovalRequestStatus = "canceled"
)

func (s ApprovalRequestStatus) IsValid() bool {
	switch s {
	case ApprovalRequestStatusPending, ApprovalRequestStatusApproved, ApprovalRequestStatusRejected, ApprovalRequestStatusCanceled:
		return true
	default:
		return false
	}
}

func (s ApprovalRequestStatus) String() string {
	return string(s)
}

func (s ApprovalRequestStatus) Parse(status string) (ApprovalRequestStatus, error) {
	requestStatus := ApprovalRequestStatus(status)
	if !requestStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidApprovalRequestStatus, status)
	}
	return requestStatus, nil
}

type ApprovalDecisionType string

const (
	ApprovalDecisionApprove ApprovalDecisionType = "approve"
	ApprovalDecisionReject  ApprovalDecisionType = "reject"
)

func (s ApprovalDecisionType) IsValid() bool {
	switch s {
	case ApprovalDecisionApprove, ApprovalDecisionReject:
		return true
	default:
		return false
	}
}

func (s ApprovalDecisionType) String() string {
	return string(s)
}

func (s ApprovalDecisionType) Parse(decision string) (ApprovalDecisionType, error) {
	decisionType := ApprovalDecisionType(decision)
	if !decisionType.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidApprovalDecision, decision)
	}
	return decisionType, nil
}

type ApprovalPolicy struct {
	ID                uuid.UUID
	Resource          string
	Action            string
	ProjectID         *uuid.UUID
	ServiceID         *uuid.UUID
	Environment       *string
	RequiredApprovals int
	Enabled           bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ApprovalRequest struct {
	ID                uuid.UUID
	Resource          string
	Action            string
	ResourceID        uuid.UUID
	RequestedBy       uuid.UUID
	ProjectID         *uuid.UUID
	ServiceID         *uuid.UUID
	Environment       *string
	Status            ApprovalRequestStatus
	RequiredApprovals int
	ApprovedCount     int
	RejectedCount     int
	ResolvedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ApprovalRequests []ApprovalRequest

type ApprovalDecision struct {
	ID                uuid.UUID
	ApprovalRequestID uuid.UUID
	DecidedBy         uuid.UUID
	Decision          ApprovalDecisionType
	Comment           string
	CreatedAt         time.Time
}
