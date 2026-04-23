package approvalrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *approvalRepositoryImpl) FindApprovalRequest(ctx context.Context, id uuid.UUID) (_ *entity.ApprovalRequest, err error) {
	const errLocation = "[repository approval/find_request FindApprovalRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalRequestsTable := table.ApprovalRequests
	stmt := postgres.SELECT(
		approvalRequestsTable.AllColumns,
	).
		FROM(approvalRequestsTable).
		WHERE(approvalRequestsTable.ID.EQ(postgres.UUID(id)))
	query, args := stmt.Sql()

	var model ApprovalRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("approval request not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying approval request", err.Error()))
	}

	request := model.toEntity()
	if request == nil {
		return nil, errs.NewInternalServerError("failed to convert approval request model to entity", nil)
	}

	return request, nil
}
