package service

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

func (r *serviceRepositoryImpl) DeleteOne(ctx context.Context, id uuid.UUID) (service *entity.Service, err error) {
	const errLocation = "[repository service/delete Delete] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	servicesTable := table.Services
	// SQL statement
	stmt := servicesTable.DELETE().
		WHERE(servicesTable.ID.EQ(postgres.UUID(id))).
		RETURNING(servicesTable.AllColumns)
	query, args := stmt.Sql()

	var model Service
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("service not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while deleting service", err.Error()))
	}

	service = model.ToEntity()
	if service == nil {
		return nil, errs.NewInternalServerError("failed to convert service model to entity", nil)
	}

	return service, nil
}
