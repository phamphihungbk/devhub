package approvalrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"github.com/google/uuid"
)

func (r *approvalRepositoryImpl) FindApprovalRequest(ctx context.Context, id uuid.UUID) (_ *entity.ApprovalRequest, err error) {
	const errLocation = "[repository approval/find_request FindApprovalRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	query := `
		SELECT
			id,
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
			rejected_count,
			resolved_at
		FROM approval_requests
		WHERE id = $1
	`

	var model approvalRequestModel
	if err := r.execer.GetContext(ctx, &model, query, id); err != nil {
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
