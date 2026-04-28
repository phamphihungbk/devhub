package service

import (
	"context"
	"encoding/json"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *serviceRepositoryImpl) CreateDependency(ctx context.Context, input *entity.ServiceDependency) (dependency *entity.ServiceDependency, err error) {
	const errLocation = "[repository service/create_dependency CreateDependency] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	config, err := json.Marshal(input.Config)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid dependency config", nil))
	}

	var port *int32
	if input.Port != nil {
		value := int32(*input.Port)
		port = &value
	}

	serviceDependenciesTable := table.ServiceDependencies
	stmt := serviceDependenciesTable.INSERT(
		serviceDependenciesTable.AllColumns.Except(serviceDependenciesTable.DefaultColumns),
	).MODEL(model.ServiceDependencies{
		ServiceID:          input.ServiceID,
		DependsOnServiceID: input.DependsOnServiceID,
		Type:               input.Type,
		Protocol:           misc.ToPointer(input.Protocol),
		Port:               port,
		Path:               misc.ToPointer(input.Path),
		Config:             string(config),
		CreatedBy:          input.CreatedBy,
	}).RETURNING(serviceDependenciesTable.AllColumns)
	query, args := stmt.Sql()

	var model ServiceDependency
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return nil, misc.WrapError(err, errs.NewConflictError("service dependency already exists", nil))
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating service dependency", err.Error()))
	}

	dependency = model.ToEntity()
	if dependency == nil {
		return nil, errs.NewInternalServerError("failed to convert service dependency model to entity", nil)
	}

	return dependency, nil
}
