package environment

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *environmentRepositoryImpl) CreateOne(ctx context.Context, input *entity.Environment) (environment *entity.Environment, err error) {
	const errLocation = "[repository environment/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	environmentsTable := table.Environments
	stmt := environmentsTable.INSERT(
		environmentsTable.ProjectID,
		environmentsTable.Name,
		environmentsTable.Tier,
		environmentsTable.Cluster,
		environmentsTable.Namespace,
		environmentsTable.ArgocdInstance,
		environmentsTable.Config,
	).MODEL(model.Environments{
		ProjectID:      input.ProjectID,
		Name:           input.Name,
		Tier:           input.Tier.String(),
		Cluster:        input.Cluster,
		Namespace:      input.Namespace,
		ArgocdInstance: input.ArgocdInstance,
		Config:         input.Config.String(),
	}).RETURNING(environmentsTable.AllColumns)

	query, args := stmt.Sql()

	var model Environment
	if err := r.execer.GetContext(ctx, &model, query, args...); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating environment", err.Error()))
	}

	environment = model.ToEntity()
	if environment == nil {
		return nil, errs.NewInternalServerError("failed to convert environment model to entity", nil)
	}

	return environment, nil
}
