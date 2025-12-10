package projectrepo

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

func (r *scaffoldRequestRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[repository scaffold_request/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	scaffoldRequestsTable := table.ScaffoldRequests
	// SQL statement
	stmt := postgres.SELECT(
		scaffoldRequestsTable.AllColumns,
	).
		FROM(table.ScaffoldRequests).
		WHERE(table.ScaffoldRequests.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model ScaffoldRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("scaffold request not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying scaffold request by id", err.Error()))
	}

	scaffoldRequest = model.ToEntity()
	if scaffoldRequest == nil {
		return nil, errs.NewInternalServerError("failed to convert scaffold request model to entity", nil)
	}

	return scaffoldRequest, nil
}
