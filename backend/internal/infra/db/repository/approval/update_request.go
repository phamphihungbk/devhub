package approvalrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *approvalRepositoryImpl) UpdateApprovalRequest(ctx context.Context, input repository.UpdateApprovalRequestInput) (_ *entity.ApprovalRequest, err error) {
	const errLocation = "[repository approval/update_request UpdateApprovalRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalRequestsTable := table.ApprovalRequests
	var updateModel ApprovalRequest
	columns := make(postgres.ColumnList, 0)

	if input.Status != nil {
		updateModel.Status = input.Status.String()
		columns = append(columns, approvalRequestsTable.Status)
	}

	if input.ApprovedCount != nil {
		updateModel.ApprovedCount = int32(misc.GetValue(input.ApprovedCount))
		columns = append(columns, approvalRequestsTable.ApprovedCount)
	}

	if input.RejectedCount != nil {
		updateModel.RejectedCount = int32(misc.GetValue(input.RejectedCount))
		columns = append(columns, approvalRequestsTable.RejectedCount)
	}

	if input.RequiredApprovals != nil {
		updateModel.RequiredApprovals = int32(misc.GetValue(input.RequiredApprovals))
		columns = append(columns, approvalRequestsTable.RequiredApprovals)
	}

	if input.ResolvedAt != nil {
		updateModel.ResolvedAt = input.ResolvedAt
		columns = append(columns, approvalRequestsTable.ResolvedAt)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	stmt := approvalRequestsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(approvalRequestsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(approvalRequestsTable.AllColumns)
	query, args := stmt.Sql()

	var model ApprovalRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("approval request not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating approval request", err.Error()))
	}

	request := model.toEntity()
	if request == nil {
		return nil, errs.NewInternalServerError("failed to convert approval request model to entity", nil)
	}

	return request, nil
}
