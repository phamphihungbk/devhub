package approvalrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *approvalRepositoryImpl) CreateApprovalPolicy(ctx context.Context, policy *entity.ApprovalPolicy) (_ *entity.ApprovalPolicy, err error) {
	const errLocation = "[repository approval/create_policy CreateApprovalPolicy] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalPoliciesTable := table.ApprovalPolicies
	// TODO: remove default value from migration
	stmt := approvalPoliciesTable.INSERT(
		approvalPoliciesTable.Resource,
		approvalPoliciesTable.ServiceID,
		approvalPoliciesTable.EnvironmentID,
		approvalPoliciesTable.RequiredApprovals,
		approvalPoliciesTable.Enabled,
	).MODEL(model.ApprovalPolicies{
		Resource:          policy.Resource,
		ServiceID:         policy.ServiceID,
		EnvironmentID:     policy.EnvironmentID,
		RequiredApprovals: int32(policy.RequiredApprovals),
		Enabled:           policy.Enabled,
	}).RETURNING(approvalPoliciesTable.AllColumns)
	query, args := stmt.Sql()

	var model ApprovalPolicy
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating approval policy", err.Error()))
	}

	created := model.toEntity()
	if created == nil {
		return nil, errs.NewInternalServerError("failed to convert approval policy model to entity", nil)
	}

	return created, nil
}
