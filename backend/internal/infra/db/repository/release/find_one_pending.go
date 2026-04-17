package releaserepo

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

func (r *releaseRepositoryImpl) FindOnePending(ctx context.Context) (release *entity.Release, err error) {
	const errLocation = "[repository release/find_one_pending FindOnePending] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	releasesTable := table.Releases
	stmt := postgres.SELECT(
		releasesTable.AllColumns,
	).
		FROM(releasesTable).
		WHERE(releasesTable.Status.EQ(postgres.String(entity.ReleaseStatusPending.String()))).
		ORDER_BY(releasesTable.CreatedAt.ASC(), releasesTable.ID.ASC()).
		LIMIT(1)

	query, args := stmt.Sql()

	var model Release
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying pending release", err.Error()))
	}

	release = model.ToEntity()
	if release == nil {
		return nil, errs.NewInternalServerError("failed to convert pending release model to entity", nil)
	}

	return release, nil
}
