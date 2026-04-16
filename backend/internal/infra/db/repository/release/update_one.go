package releaserepo

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

func (r *releaseRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateReleaseInput) (release *entity.Release, err error) {
	const errLocation = "[repository release/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	releasesTable := table.Releases
	var updateModel Release
	columns := make(postgres.ColumnList, 0)

	if input.Status != nil {
		updateModel.Status = string(*input.Status)
		columns = append(columns, releasesTable.Status)
	}
	if input.ExternalRef != nil {
		updateModel.ExternalRef = *input.ExternalRef
		columns = append(columns, releasesTable.ExternalRef)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	stmt := releasesTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(releasesTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(releasesTable.AllColumns)
	query, args := stmt.Sql()

	var model Release
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("release not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating release", err.Error()))
	}

	release = model.ToEntity()
	if release == nil {
		return nil, errs.NewInternalServerError("failed to convert release model to entity", nil)
	}

	return release, nil
}
