package usecase

import (
	"context"
	"slices"
	"strings"
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
	ActorName string
	Summary   string
	Comment   string
	CreatedAt time.Time
}

type FindApprovalRequestDetailOutput struct {
	ApprovalRequest *entity.ApprovalRequest
	ResourceName    string
	RequestedByName string
	Scope           string
	ActorNames      map[uuid.UUID]string
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

	resourceName := entity.ApprovalResource(request.Resource).String()
	actorNames := u.loadApprovalActorNames(ctx, request, decisions)
	requestedByName := actorNames[request.RequestedBy]

	return &FindApprovalRequestDetailOutput{
		ApprovalRequest: request,
		ResourceName:    resourceName,
		RequestedByName: requestedByName,
		ActorNames:      actorNames,
		Decisions:       decisions,
		AuditEvents:     newApprovalAuditEvents(request, decisions, actorNames),
	}, nil
}

func newApprovalAuditEvents(request *entity.ApprovalRequest, decisions []entity.ApprovalDecision, actorNames map[uuid.UUID]string) []ApprovalAuditEvent {
	if request == nil {
		return []ApprovalAuditEvent{}
	}

	events := []ApprovalAuditEvent{
		{
			Type:      "requested",
			ActorID:   request.RequestedBy,
			ActorName: actorNames[request.RequestedBy],
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
			ActorName: actorNames[decision.DecidedBy],
			Summary:   summary,
			Comment:   decision.Comment,
			CreatedAt: decision.CreatedAt,
		})
	}

	if request.Status != entity.ApprovalRequestStatusPending && !request.ResolvedAt.IsZero() {
		events = append(events, ApprovalAuditEvent{
			Type:      "resolved",
			ActorID:   request.RequestedBy,
			ActorName: actorNames[request.RequestedBy],
			Summary:   "Approval request resolved as " + request.Status.String(),
			CreatedAt: request.ResolvedAt,
		})
	}

	slices.SortStableFunc(events, func(a, b ApprovalAuditEvent) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	return events
}

func (u *approvalUsecase) loadApprovalActorNames(ctx context.Context, request *entity.ApprovalRequest, decisions []entity.ApprovalDecision) map[uuid.UUID]string {
	actorNames := map[uuid.UUID]string{}
	if request == nil {
		return actorNames
	}

	actorIDs := []uuid.UUID{request.RequestedBy}
	for _, decision := range decisions {
		actorIDs = append(actorIDs, decision.DecidedBy)
	}

	for _, actorID := range actorIDs {
		if _, exists := actorNames[actorID]; exists {
			continue
		}

		user, err := u.userRepository.FindOne(ctx, actorID)
		if err != nil || user == nil {
			actorNames[actorID] = actorID.String()
			continue
		}

		if strings.TrimSpace(user.Name) != "" {
			actorNames[actorID] = user.Name
			continue
		}
		actorNames[actorID] = user.Email
	}

	return actorNames
}

func (u *approvalUsecase) environmentName(ctx context.Context, id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}

	environment, err := u.environmentRepository.FindOne(ctx, id)
	if err != nil || environment == nil || strings.TrimSpace(environment.Name) == "" {
		return id.String()
	}

	return environment.Name
}
