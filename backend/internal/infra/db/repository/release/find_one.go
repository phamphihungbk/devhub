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
	"github.com/google/uuid"
)

func (r *releaseRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (release *entity.Release, err error) {
	const errLocation = "[repository release/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	releasesTable := table.Releases
	stmt := postgres.SELECT(
		releasesTable.AllColumns,
	).
		FROM(releasesTable).
		WHERE(releasesTable.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model Release
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("release not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying release by id", err.Error()))
	}

	release = model.ToEntity()
	if release == nil {
		return nil, errs.NewInternalServerError("failed to convert release model to entity", nil)
	}

	return release, nil
}
