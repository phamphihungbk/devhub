package usecase

import (
	"context"
	"slices"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type FindApprovalRequestDetailInput struct {
	ID string `json:"id" validate:"required,uuid"`
}

type ApprovalAuditEvent struct {
	Type      string
	ActorID   uuid.UUID
	Summary   string
	Comment   string
	CreatedAt time.Time
}

type FindApprovalRequestDetailOutput struct {
	ApprovalRequest *entity.ApprovalRequest
	Decisions       []entity.ApprovalDecision
	AuditEvents     []ApprovalAuditEvent
}

func (u *approvalUsecase) FindApprovalRequestDetail(ctx context.Context, input FindApprovalRequestDetailInput) (_ *FindApprovalRequestDetailOutput, err error) {
	const errLocation = "[usecase approval/find_request_detail FindApprovalRequestDetail] "
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

	approvalRequestID := uuid.MustParse(input.ID)
	request, err := u.approvalRepository.FindApprovalRequest(ctx, approvalRequestID)
	if err != nil {
		return nil, err
	}

	decisions, err := u.approvalRepository.FindApprovalDecisions(ctx, approvalRequestID)
	if err != nil {
		return nil, err
	}

	return &FindApprovalRequestDetailOutput{
		ApprovalRequest: request,
		Decisions:       decisions,
		AuditEvents:     newApprovalAuditEvents(request, decisions),
	}, nil
}

func newApprovalAuditEvents(request *entity.ApprovalRequest, decisions []entity.ApprovalDecision) []ApprovalAuditEvent {
	if request == nil {
		return []ApprovalAuditEvent{}
	}

	events := []ApprovalAuditEvent{
		{
			Type:      "requested",
			ActorID:   request.RequestedBy,
			Summary:   "Approval request created",
			CreatedAt: request.CreatedAt,
		},
	}

	for _, decision := range decisions {
		summary := "Request approved"
		if decision.Decision == entity.ApprovalDecisionReject {
			summary = "Request rejected"
		}
		events = append(events, ApprovalAuditEvent{
			Type:      decision.Decision.String(),
			ActorID:   decision.DecidedBy,
			Summary:   summary,
			Comment:   decision.Comment,
			CreatedAt: decision.CreatedAt,
		})
	}

	if request.ResolvedAt != nil {
		events = append(events, ApprovalAuditEvent{
			Type:      "resolved",
			ActorID:   request.RequestedBy,
			Summary:   "Approval request resolved as " + request.Status.String(),
			CreatedAt: *request.ResolvedAt,
		})
	}

	slices.SortStableFunc(events, func(a, b ApprovalAuditEvent) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	return events
}
