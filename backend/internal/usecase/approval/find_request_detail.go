package usecase

import (
	"context"
	"fmt"
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

	resourceName, scope := u.describeApprovalTarget(ctx, request)
	actorNames := u.loadApprovalActorNames(ctx, request, decisions)
	requestedByName := actorNames[request.RequestedBy]

	return &FindApprovalRequestDetailOutput{
		ApprovalRequest: request,
		ResourceName:    resourceName,
		RequestedByName: requestedByName,
		Scope:           scope,
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

	if request.ResolvedAt != nil {
		events = append(events, ApprovalAuditEvent{
			Type:      "resolved",
			ActorID:   request.RequestedBy,
			ActorName: actorNames[request.RequestedBy],
			Summary:   "Approval request resolved as " + request.Status.String(),
			CreatedAt: *request.ResolvedAt,
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

func (u *approvalUsecase) describeApprovalTarget(ctx context.Context, request *entity.ApprovalRequest) (string, string) {
	if request == nil {
		return "", "Global"
	}

	resourceName := request.ResourceID.String()
	scope := u.formatApprovalScope(ctx, request.ProjectID, request.ServiceID, request.Environment)

	switch entity.ApprovalResource(request.Resource) {
	case entity.ApprovalResourceScaffoldRequest:
		scaffoldRequest, err := u.scaffoldRequestRepository.FindOne(ctx, request.ResourceID)
		if err == nil && scaffoldRequest != nil {
			if strings.TrimSpace(scaffoldRequest.Variables.ServiceName) != "" {
				resourceName = scaffoldRequest.Variables.ServiceName
			}
			if scope == "Global" {
				env := scaffoldRequest.Environment.String()
				scope = u.formatApprovalScope(ctx, &scaffoldRequest.ProjectID, nil, &env)
			}
		}
	case entity.ApprovalResourceDeployment:
		deployment, err := u.deploymentRepository.FindOne(ctx, request.ResourceID)
		if err == nil && deployment != nil {
			serviceName := deployment.ServiceID.String()
			if service, err := u.serviceRepository.FindOne(ctx, deployment.ServiceID); err == nil && service != nil {
				serviceName = service.Name
			}
			resourceName = fmt.Sprintf("%s %s to %s", serviceName, deployment.Version, deployment.Environment.String())
			if scope == "Global" {
				env := deployment.Environment.String()
				scope = u.formatApprovalScope(ctx, nil, &deployment.ServiceID, &env)
			}
		}
	case entity.ApprovalResourceRelease:
		release, err := u.releaseRepository.FindOne(ctx, request.ResourceID)
		if err == nil && release != nil {
			serviceName := release.ServiceID.String()
			if service, err := u.serviceRepository.FindOne(ctx, release.ServiceID); err == nil && service != nil {
				serviceName = service.Name
			}
			name := release.Tag
			if strings.TrimSpace(release.Name) != "" {
				name = release.Name
			}
			resourceName = fmt.Sprintf("%s %s", serviceName, name)
			if scope == "Global" {
				scope = u.formatApprovalScope(ctx, nil, &release.ServiceID, nil)
			}
		}
	}

	return resourceName, scope
}

func (u *approvalUsecase) formatApprovalScope(ctx context.Context, projectID *uuid.UUID, serviceID *uuid.UUID, environment *string) string {
	parts := make([]string, 0, 3)

	if projectID != nil {
		projectName := projectID.String()
		if project, err := u.projectRepository.FindOne(ctx, *projectID); err == nil && project != nil {
			projectName = project.Name
		}
		parts = append(parts, projectName)
	}

	if serviceID != nil {
		serviceName := serviceID.String()
		if service, err := u.serviceRepository.FindOne(ctx, *serviceID); err == nil && service != nil {
			serviceName = service.Name
			if projectID == nil {
				projectName := service.ProjectID.String()
				if project, err := u.projectRepository.FindOne(ctx, service.ProjectID); err == nil && project != nil {
					projectName = project.Name
				}
				parts = append(parts, projectName)
			}
		}
		parts = append(parts, serviceName)
	}

	if environment != nil && strings.TrimSpace(*environment) != "" {
		parts = append(parts, strings.TrimSpace(*environment))
	}

	if len(parts) == 0 {
		return "Global"
	}

	return strings.Join(parts, " / ")
}
