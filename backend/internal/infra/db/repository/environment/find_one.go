package environment

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

func (r *environmentRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (environment *entity.Environment, err error) {
	const errLocation = "[repository environment/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	environmentsTable := table.Environments
	stmt := postgres.SELECT(
		environmentsTable.AllColumns,
	).
		FROM(environmentsTable).
		WHERE(environmentsTable.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model Environment
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("environment not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying environment by id", err.Error()))
	}

	environment = model.ToEntity()
	if environment == nil {
		return nil, errs.NewInternalServerError("failed to convert environment model to entity", nil)
	}

	return environment, nil
}
