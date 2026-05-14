package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidApprovalRequestStatus = fmt.Errorf("invalid approval request status")
)

type ApprovalRequestStatus string

const (
	ApprovalRequestStatusPending  ApprovalRequestStatus = "pending"
	ApprovalRequestStatusApproved ApprovalRequestStatus = "approved"
	ApprovalRequestStatusRejected ApprovalRequestStatus = "rejected"
)

func (s ApprovalRequestStatus) IsValid() bool {
	switch s {
	case ApprovalRequestStatusPending, ApprovalRequestStatusApproved, ApprovalRequestStatusRejected:
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

type ApprovalRequest struct {
	ID                uuid.UUID
	Resource          string
	Action            string
	ResourceID        uuid.UUID
	RequestedBy       uuid.UUID
	Status            ApprovalRequestStatus
	RequiredApprovals int
	ApprovedCount     int
	RejectedCount     int
	ResolvedAt        time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ApprovalRequests []ApprovalRequest
