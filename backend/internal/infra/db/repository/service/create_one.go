package service

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *serviceRepositoryImpl) CreateOne(ctx context.Context, input *entity.Service) (service *entity.Service, err error) {
	const errLocation = "[repository service/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	servicesTable := table.Services
	stmt := servicesTable.INSERT(
		servicesTable.AllColumns.Except(servicesTable.DefaultColumns),
	).MODEL(model.Services{
		ProjectID: input.ProjectID,
		Name:      input.Name,
		RepoURL:   input.RepoURL,
	}).RETURNING(servicesTable.AllColumns)

	query, args := stmt.Sql()

	var model Service
	if err := r.execer.GetContext(
		ctx,
		&model,
		query,
		args...,
	); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating service", err.Error()))
	}

	service = model.ToEntity()
	if service == nil {
		return nil, errs.NewInternalServerError("failed to convert service model to entity", nil)
	}

	return service, nil
}
