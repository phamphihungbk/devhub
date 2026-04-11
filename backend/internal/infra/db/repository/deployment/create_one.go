package deploymentrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *deploymentRepositoryImpl) CreateOne(ctx context.Context, input *entity.Deployment) (deployment *entity.Deployment, err error) {
	const errLocation = "[repository deployment/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	deploymentsTable := table.Deployments
	// SQL statement
	stmt := deploymentsTable.INSERT(
		deploymentsTable.AllColumns.Except(deploymentsTable.DefaultColumns), // Exclude columns with default values
	).MODEL(model.Deployments{
		ProjectID:   input.ProjectID,
		Environment: input.Environment.String(),
		Service:     input.Service,
		Version:     input.Version,
		Status:      input.Status.String(),
		TriggeredBy: input.TriggeredBy,
	}).RETURNING(deploymentsTable.AllColumns)
	query, args := stmt.Sql()

	var model Deployment
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating deployment", err.Error()))
	}

	deployment = model.ToEntity()
	if deployment == nil {
		return nil, errs.NewInternalServerError("failed to convert deployment model to entity", nil)
	}

	return deployment, nil
}
