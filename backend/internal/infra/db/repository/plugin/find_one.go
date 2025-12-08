package userrepo

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

func (r *pluginRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (plugin *entity.Plugin, err error) {
	const errLocation = "[repository plugin/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	pluginsTable := table.Plugins
	// SQL statement
	stmt := postgres.SELECT(
		pluginsTable.AllColumns,
	).
		FROM(pluginsTable).
		WHERE(pluginsTable.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model Plugin
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("plugin not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying plugin by id", err.Error()))
	}

	plugin = model.ToEntity()
	if plugin == nil {
		return nil, errs.NewInternalServerError("failed to convert plugin model to entity", nil)
	}

	return plugin, nil
}
