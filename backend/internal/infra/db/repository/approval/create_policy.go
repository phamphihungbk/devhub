package approvalrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *approvalRepositoryImpl) CreateApprovalPolicy(ctx context.Context, policy *entity.ApprovalPolicy) (_ *entity.ApprovalPolicy, err error) {
	const errLocation = "[repository approval/create_policy CreateApprovalPolicy] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	query := `
		INSERT INTO approval_policies (
			resource,
			action,
			project_id,
			service_id,
			environment,
			required_approvals,
			enabled
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, resource, action, project_id, service_id, environment, required_approvals, enabled, created_at, updated_at
	`

	var model approvalPolicyModel
	if err := r.execer.GetContext(
		ctx,
		&model,
		query,
		policy.Resource,
		policy.Action,
		policy.ProjectID,
		policy.ServiceID,
		policy.Environment,
		policy.RequiredApprovals,
		policy.Enabled,
	); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating approval policy", err.Error()))
	}

	created := model.toEntity()
	if created == nil {
		return nil, errs.NewInternalServerError("failed to convert approval policy model to entity", nil)
	}

	return created, nil
}
