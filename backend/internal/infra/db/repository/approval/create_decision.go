package approvalrepo

import (
	"context"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *approvalRepositoryImpl) CreateApprovalDecision(ctx context.Context, decision *entity.ApprovalDecision) (_ *entity.ApprovalDecision, err error) {
	const errLocation = "[repository approval/create_decision CreateApprovalDecision] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	query := `
		INSERT INTO approval_decisions (
			approval_request_id,
			decided_by,
			decision,
			comment
		) VALUES ($1, $2, $3, $4)
		RETURNING id, approval_request_id, decided_by, decision, comment
	`

	var model approvalDecisionModel
	if err := r.execer.GetContext(
		ctx,
		&model,
		query,
		decision.ApprovalRequestID,
		decision.DecidedBy,
		decision.Decision,
		decision.Comment,
	); err != nil {
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
