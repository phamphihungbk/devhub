package scaffoldrequestrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *scaffoldRequestRepositoryImpl) FindOnePending(ctx context.Context) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[repository scaffold_request/find_one_pending FindOnePending] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	scaffoldRequestsTable := table.ScaffoldRequests
	stmt := postgres.SELECT(
		scaffoldRequestsTable.AllColumns,
	).
		FROM(scaffoldRequestsTable).
		WHERE(scaffoldRequestsTable.Status.EQ(postgres.String(entity.ScaffoldRequestPending.String()))).
		ORDER_BY(scaffoldRequestsTable.CreatedAt.ASC(), scaffoldRequestsTable.ID.ASC()).
		LIMIT(1)

	query, args := stmt.Sql()

	var model ScaffoldRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying pending scaffold request", err.Error()))
	}

	scaffoldRequest = model.ToEntity()
	if scaffoldRequest == nil {
		return nil, errs.NewInternalServerError("failed to convert pending scaffold request model to entity", nil)
	}

	return scaffoldRequest, nil
}
