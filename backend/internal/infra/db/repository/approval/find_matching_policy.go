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

func (r *approvalRepositoryImpl) FindMatchingApprovalPolicy(ctx context.Context, input repository.FindMatchingApprovalPolicyInput) (_ *entity.ApprovalPolicy, err error) {
	const errLocation = "[repository approval/find_matching_policy FindMatchingApprovalPolicy] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalPoliciesTable := table.ApprovalPolicies
	stmt := postgres.SELECT(
		approvalPoliciesTable.AllColumns,
	).FROM(approvalPoliciesTable).
		WHERE(
			approvalPoliciesTable.Resource.EQ(postgres.String(input.Resource)).
				AND(approvalPoliciesTable.Action.EQ(postgres.String(input.Action))).
				AND(approvalPoliciesTable.Enabled.EQ(postgres.Bool(true))).
				AND(
					approvalPoliciesTable.ProjectID.IS_NULL().
						OR(approvalPoliciesTable.ProjectID.EQ(postgres.UUID(misc.GetValue(input.ProjectID)))),
				).
				AND(
					approvalPoliciesTable.ServiceID.IS_NULL().
						OR(approvalPoliciesTable.ServiceID.EQ(postgres.UUID(misc.GetValue(input.ServiceID)))),
				).
				AND(
					approvalPoliciesTable.Environment.IS_NULL().
						OR(approvalPoliciesTable.Environment.EQ(postgres.String(misc.GetValue(input.Environment)))),
				),
		).
		ORDER_BY(
			postgres.Raw("CASE WHEN service_id IS NULL THEN 0 ELSE 1 END DESC"),
			postgres.Raw("CASE WHEN project_id IS NULL THEN 0 ELSE 1 END DESC"),
			postgres.Raw("CASE WHEN environment IS NULL THEN 0 ELSE 1 END DESC"),
			approvalPoliciesTable.UpdatedAt.DESC(),
			approvalPoliciesTable.CreatedAt.DESC(),
		).
		LIMIT(1)
	query, args := stmt.Sql()

	var model ApprovalPolicy
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
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
