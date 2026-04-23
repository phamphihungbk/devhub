package usecase

import (
	"context"
	"strings"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateApprovalDecisionInput struct {
	ApprovalRequestID string `json:"approval_request_id" validate:"required,uuid"`
	DecidedBy         string `json:"decided_by" validate:"required,uuid"`
	Decision          string `json:"decision" validate:"required,oneof=approve reject"`
	Comment           string `json:"comment"`
}

type CreateApprovalDecisionOutput struct {
	ApprovalRequest *entity.ApprovalRequest
	Decision        *entity.ApprovalDecision
}

func (u *approvalUsecase) CreateApprovalDecision(ctx context.Context, input CreateApprovalDecisionInput) (_ *CreateApprovalDecisionOutput, err error) {
	const errLocation = "[usecase approval/create_decision CreateApprovalDecision] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err := vInstance.Struct(input); err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	approvalRequestID := uuid.MustParse(input.ApprovalRequestID)
	decidedBy := uuid.MustParse(input.DecidedBy)
	decisionType, err := new(entity.ApprovalDecisionType).Parse(input.Decision)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid approval decision", nil))
	}

	request, err := u.approvalRepository.FindApprovalRequest(ctx, approvalRequestID)
	if err != nil {
		return nil, err
	}
	if request.Status != entity.ApprovalRequestStatusPending {
		return nil, errs.NewConflictError("approval request is already resolved", map[string]any{"status": request.Status.String()})
	}

	decision, err := u.approvalRepository.CreateApprovalDecision(ctx, &entity.ApprovalDecision{
		ApprovalRequestID: approvalRequestID,
		DecidedBy:         decidedBy,
		Decision:          decisionType,
		Comment:           strings.TrimSpace(input.Comment),
	})
	if err != nil {
		return nil, err
	}

	approvedCount := request.ApprovedCount
	rejectedCount := request.RejectedCount
	status := request.Status
	var resolvedAt *time.Time

	switch decisionType {
	case entity.ApprovalDecisionApprove:
		approvedCount++
		if approvedCount >= request.RequiredApprovals {
			status = entity.ApprovalRequestStatusApproved
			now := time.Now()
			resolvedAt = &now
		}
	case entity.ApprovalDecisionReject:
		rejectedCount++
		status = entity.ApprovalRequestStatusRejected
		now := time.Now()
		resolvedAt = &now
	}

	updatedRequest, err := u.approvalRepository.UpdateApprovalRequest(ctx, repository.UpdateApprovalRequestInput{
		ID:            approvalRequestID,
		Status:        &status,
		ApprovedCount: &approvedCount,
		RejectedCount: &rejectedCount,
		ResolvedAt:    resolvedAt,
	})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update approval request", nil))
	}

	if decisionType == entity.ApprovalDecisionApprove && updatedRequest.Status == entity.ApprovalRequestStatusApproved && updatedRequest.Resource == "scaffold_request" {
		now := misc.GetValue(resolvedAt)
		scaffoldStatus := entity.ScaffoldRequestApproved
		if _, err := u.scaffoldRequestRepository.UpdateOne(ctx, repository.UpdateScaffoldRequestInput{
			ID:         updatedRequest.ResourceID,
			Status:     &scaffoldStatus,
			ApprovedBy: &decidedBy,
			ApprovedAt: &now,
		}); err != nil {
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update scaffold request after approval", nil))
		}
	}

	return &CreateApprovalDecisionOutput{
		ApprovalRequest: updatedRequest,
		Decision:        decision,
	}, nil
}
