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

func (r *deploymentRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (deployment *entity.Deployment, err error) {
	const errLocation = "[repository deployment/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	deploymentsTable := table.Deployments
	// SQL statement
	stmt := postgres.SELECT(
		deploymentsTable.AllColumns,
	).
		FROM(table.Deployments).
		WHERE(table.Deployments.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model Deployment
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("deployment not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying deployment by id", err.Error()))
	}

	deployment = model.ToEntity()
	if deployment == nil {
		return nil, errs.NewInternalServerError("failed to convert deployment model to entity", nil)
	}

	return deployment, nil
}
