package pluginrepo

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

func (r *pluginRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdatePluginInput) (plugin *entity.Plugin, err error) {
	const errLocation = "[repository plugin/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	pluginsTable := table.Plugins
	var updateModel Plugin
	columns := make(postgres.ColumnList, 0)

	// build the update model
	if input.Name != nil {
		updateModel.Name = misc.GetValue(input.Name)
		columns = append(columns, pluginsTable.Name)
	}
	if input.Description != nil {
		updateModel.Description = input.Description
		columns = append(columns, pluginsTable.Description)
	}
	if input.Type != nil {
		updateModel.Type = string(*input.Type)
		columns = append(columns, pluginsTable.Type)
	}
	if input.Version != nil {
		updateModel.Version = misc.GetValue(input.Version)
		columns = append(columns, pluginsTable.Version)
	}
	if input.Entrypoint != nil {
		updateModel.Entrypoint = misc.GetValue(input.Entrypoint)
		columns = append(columns, pluginsTable.Entrypoint)
	}
	if input.Scope != nil {
		updateModel.Scope = misc.GetValue(input.Scope)
		columns = append(columns, pluginsTable.Scope)
	}
	if input.Enabled != nil {
		updateModel.Enabled = misc.GetValue(input.Enabled)
		columns = append(columns, pluginsTable.Enabled)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	// SQL statement
	stmt := pluginsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(pluginsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(pluginsTable.AllColumns)
	query, args := stmt.Sql()

	var model Plugin
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("plugin not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating plugin", err.Error()))
	}

	plugin = model.ToEntity()
	if plugin == nil {
		return nil, errs.NewInternalServerError("failed to convert plugin model to entity", nil)
	}

	return plugin, nil
}
