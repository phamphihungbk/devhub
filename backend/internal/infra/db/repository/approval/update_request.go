package approvalrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
)

func (r *approvalRepositoryImpl) UpdateApprovalRequest(ctx context.Context, input repository.UpdateApprovalRequestInput) (_ *entity.ApprovalRequest, err error) {
	const errLocation = "[repository approval/update_request UpdateApprovalRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	query := `
		UPDATE approval_requests
		SET
			status = COALESCE($2, status),
			approved_count = COALESCE($3, approved_count),
			rejected_count = COALESCE($4, rejected_count),
			required_approvals = COALESCE($5, required_approvals),
			resolved_at = CASE
				WHEN $6 IS NULL THEN resolved_at
				ELSE $6
			END,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, resource, action, resource_id, requested_by, project_id, service_id, environment, status, required_approvals, approved_count, rejected_count, resolved_at
	`

	var (
		statusValue *string
	)

	if input.Status != nil {
		value := input.Status.String()
		statusValue = &value
	}

	var model approvalRequestModel
	if err := r.execer.GetContext(
		ctx,
		&model,
		query,
		input.ID,
		statusValue,
		input.ApprovedCount,
		input.RejectedCount,
		input.RequiredApprovals,
		input.ResolvedAt,
	); err != nil {
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
