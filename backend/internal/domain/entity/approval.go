package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidApprovalRequestStatus = fmt.Errorf("invalid approval request status")
	ErrInvalidApprovalDecision      = fmt.Errorf("invalid approval decision")
)

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

type ApprovalDecision struct {
	ID                uuid.UUID
	ApprovalRequestID uuid.UUID
	DecidedBy         uuid.UUID
	Decision          ApprovalDecisionType
	Comment           string
	CreatedAt         time.Time
}
