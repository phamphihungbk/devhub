package scaffoldrequestrepo

import (
	"context"
	"database/sql"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
	"errors"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *scaffoldRequestRepositoryImpl) DeleteOne(ctx context.Context, id uuid.UUID) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[repository scaffold_request/delete Delete] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	scaffoldRequestsTable := table.ScaffoldRequests
	// SQL statement
	stmt := scaffoldRequestsTable.DELETE().
		WHERE(scaffoldRequestsTable.ID.EQ(postgres.UUID(id))).
		RETURNING(scaffoldRequestsTable.AllColumns)
	query, args := stmt.Sql()

	var model ScaffoldRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("scaffold request not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while deleting scaffold request", err.Error()))
	}

	scaffoldRequest = model.ToEntity()
	if scaffoldRequest == nil {
		return nil, errs.NewInternalServerError("failed to convert scaffoldRequest model to entity", nil)
	}

	return scaffoldRequest, nil
}
