package approvalrepo

import (
	"context"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *approvalRepositoryImpl) CreateApprovalDecision(ctx context.Context, decision *entity.ApprovalDecision) (_ *entity.ApprovalDecision, err error) {
	const errLocation = "[repository approval/create_decision CreateApprovalDecision] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalDecisionsTable := table.ApprovalDecisions
	stmt := approvalDecisionsTable.INSERT(
		approvalDecisionsTable.AllColumns.Except(approvalDecisionsTable.DefaultColumns),
	).MODEL(model.ApprovalDecisions{
		ApprovalRequestID: decision.ApprovalRequestID,
		DecidedBy:         decision.DecidedBy,
		Decision:          decision.Decision.String(),
		Comment:           misc.ToPointer(decision.Comment),
	}).RETURNING(approvalDecisionsTable.AllColumns)
	query, args := stmt.Sql()

	var model ApprovalDecision
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return nil, misc.WrapError(err, errs.NewConflictError("approval decision already exists for this user", nil))
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating approval decision", err.Error()))
	}

	created := model.toEntity()
	if created == nil {
		return nil, errs.NewInternalServerError("failed to convert approval decision model to entity", nil)
	}

	return created, nil
}
