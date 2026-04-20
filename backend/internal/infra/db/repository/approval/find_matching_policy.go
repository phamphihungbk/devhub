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

func (r *approvalRepositoryImpl) FindMatchingApprovalPolicy(ctx context.Context, input repository.FindMatchingApprovalPolicyInput) (_ *entity.ApprovalPolicy, err error) {
	const errLocation = "[repository approval/find_matching_policy FindMatchingApprovalPolicy] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	query := `
		SELECT
			id,
			resource,
			action,
			project_id,
			service_id,
			environment,
			required_approvals,
			enabled,
			created_at,
			updated_at
		FROM approval_policies
		WHERE resource = $1
			AND action = $2
			AND enabled = TRUE
			AND (project_id IS NULL OR project_id = $3)
			AND (service_id IS NULL OR service_id = $4)
			AND (environment IS NULL OR environment = $5)
		ORDER BY
			CASE WHEN service_id IS NULL THEN 0 ELSE 1 END DESC,
			CASE WHEN project_id IS NULL THEN 0 ELSE 1 END DESC,
			CASE WHEN environment IS NULL THEN 0 ELSE 1 END DESC,
			updated_at DESC,
			created_at DESC
		LIMIT 1
	`

	var model approvalPolicyModel
	if err := r.execer.GetContext(
		ctx,
		&model,
		query,
		input.Resource,
		input.Action,
		input.ProjectID,
		input.ServiceID,
		input.Environment,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying approval policy", err.Error()))
	}

	policy := model.toEntity()
	if policy == nil {
		return nil, errs.NewInternalServerError("failed to convert approval policy model to entity", nil)
	}

	return policy, nil
}
