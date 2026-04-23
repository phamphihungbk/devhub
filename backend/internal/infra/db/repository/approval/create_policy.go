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
	stmt := approvalPoliciesTable.INSERT(
		approvalPoliciesTable.AllColumns.Except(approvalPoliciesTable.DefaultColumns),
	).MODEL(model.ApprovalPolicies{
		Resource:          policy.Resource,
		Action:            policy.Action,
		ProjectID:         policy.ProjectID,
		ServiceID:         policy.ServiceID,
		Environment:       policy.Environment,
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
