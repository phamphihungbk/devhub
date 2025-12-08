package userrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *pluginRepositoryImpl) CreateOne(ctx context.Context, input *entity.Plugin) (plugin *entity.Plugin, err error) {
	const errLocation = "[repository plugin/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	pluginsTable := table.Plugins
	// SQL statement
	stmt := pluginsTable.INSERT(
		pluginsTable.AllColumns.Except(pluginsTable.DefaultColumns), // Exclude columns with default values
	).MODEL(model.Plugins{
		Name:  input.Name,
		Email: input.Email,
		Role:  string(input.Role),
	}).RETURNING(pluginsTable.AllColumns)
	query, args := stmt.Sql()

	var model Plugin
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating plugin", err.Error()))
	}

	plugin = model.ToEntity()
	if plugin == nil {
		return nil, errs.NewInternalServerError("failed to convert plugin model to entity", nil)
	}

	return plugin, nil
}
