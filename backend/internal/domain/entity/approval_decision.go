package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidApprovalDecision = fmt.Errorf("invalid approval decision")
)

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

type ApprovalDecision struct {
	ID                uuid.UUID
	ApprovalRequestID uuid.UUID
	DecidedBy         uuid.UUID
	Decision          ApprovalDecisionType
	Comment           string
	CreatedAt         time.Time
}

type ApprovalDecisions []ApprovalDecision
