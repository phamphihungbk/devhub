package service

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

func (r *serviceRepositoryImpl) DeleteDependency(ctx context.Context, serviceID uuid.UUID, dependencyID uuid.UUID) (dependency *entity.ServiceDependency, err error) {
	const errLocation = "[repository service/delete_dependency DeleteDependency] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	serviceDependenciesTable := table.ServiceDependencies
	stmt := serviceDependenciesTable.DELETE().
		WHERE(
			serviceDependenciesTable.ID.EQ(postgres.UUID(dependencyID)).
				AND(serviceDependenciesTable.ServiceID.EQ(postgres.UUID(serviceID))),
		).
		RETURNING(serviceDependenciesTable.AllColumns)
	query, args := stmt.Sql()

	var model ServiceDependency
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("service dependency not found", nil)
		}

		return nil, misc.WrapError(err, errs.NewDatabaseError("error while deleting service dependency", err.Error()))
	}

	dependency = model.ToEntity()
	if dependency == nil {
		return nil, errs.NewInternalServerError("failed to convert service dependency model to entity", nil)
	}

	return dependency, nil
}
