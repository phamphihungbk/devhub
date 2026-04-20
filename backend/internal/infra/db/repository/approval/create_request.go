package approvalrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *approvalRepositoryImpl) CreateApprovalRequest(ctx context.Context, request *entity.ApprovalRequest) (_ *entity.ApprovalRequest, err error) {
	const errLocation = "[repository approval/create_request CreateApprovalRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	query := `
		INSERT INTO approval_requests (
			resource,
			action,
			resource_id,
			requested_by,
			project_id,
			service_id,
			environment,
			status,
			required_approvals,
			approved_count,
			rejected_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, resource, action, resource_id, requested_by, project_id, service_id, environment, status, required_approvals, approved_count, rejected_count, resolved_at
	`

	var model approvalRequestModel
	if err := r.execer.GetContext(
		ctx,
		&model,
		query,
		request.Resource,
		request.Action,
		request.ResourceID,
		request.RequestedBy,
		request.ProjectID,
		request.ServiceID,
		request.Environment,
		request.Status,
		request.RequiredApprovals,
		request.ApprovedCount,
		request.RejectedCount,
	); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating approval request", err.Error()))
	}

	created := model.toEntity()
	if created == nil {
		return nil, errs.NewInternalServerError("failed to convert approval request model to entity", nil)
	}

	return created, nil
}
