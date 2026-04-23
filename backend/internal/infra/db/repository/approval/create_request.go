package approvalrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *approvalRepositoryImpl) CreateApprovalRequest(ctx context.Context, request *entity.ApprovalRequest) (_ *entity.ApprovalRequest, err error) {
	const errLocation = "[repository approval/create_request CreateApprovalRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalRequestsTable := table.ApprovalRequests
	stmt := approvalRequestsTable.INSERT(
		approvalRequestsTable.AllColumns.Except(approvalRequestsTable.DefaultColumns),
	).MODEL(model.ApprovalRequests{
		Resource:          request.Resource,
		Action:            request.Action,
		ResourceID:        request.ResourceID,
		RequestedBy:       request.RequestedBy,
		ProjectID:         request.ProjectID,
		ServiceID:         request.ServiceID,
		Environment:       request.Environment,
		Status:            request.Status.String(),
		RequiredApprovals: int32(request.RequiredApprovals),
		ApprovedCount:     int32(request.ApprovedCount),
		RejectedCount:     int32(request.RejectedCount),
	}).RETURNING(approvalRequestsTable.AllColumns)
	query, args := stmt.Sql()

	var model ApprovalRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating approval request", err.Error()))
	}

	created := model.toEntity()
	if created == nil {
		return nil, errs.NewInternalServerError("failed to convert approval request model to entity", nil)
	}

	return created, nil
}
