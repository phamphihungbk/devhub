package approvalrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *approvalRepositoryImpl) FindApprovalDecisions(ctx context.Context, approvalRequestID uuid.UUID) (_ []entity.ApprovalDecision, err error) {
	const errLocation = "[repository approval/find_decisions FindApprovalDecisions] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalDecisionsTable := table.ApprovalDecisions
	stmt := postgres.SELECT(
		approvalDecisionsTable.AllColumns,
	).
		FROM(approvalDecisionsTable).
		WHERE(approvalDecisionsTable.ApprovalRequestID.EQ(postgres.UUID(approvalRequestID))).
		ORDER_BY(approvalDecisionsTable.CreatedAt.ASC())
	query, args := stmt.Sql()

	var models ApprovalDecisions
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying approval decisions", err.Error()))
	}

	return models.ToEntities(), nil
}
