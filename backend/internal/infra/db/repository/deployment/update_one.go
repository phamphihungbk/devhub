package userrepo

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

func (r *deploymentRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateUserInput) (deployment *entity.Deployment, err error) {
	const errLocation = "[repository deployment/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	deploymentsTable := table.Deployments
	var updateModel Deployment
	columns := make(postgres.ColumnList, 0)

	// build the update model
	if input.Name != nil {
		updateModel.Name = string(*input.Name)
		columns = append(columns, usersTable.Name)
	}
	if input.Role != nil {
		updateModel.Role = string(*input.Role)
		columns = append(columns, deploymentsTable.Role)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	// SQL statement
	stmt := deploymentsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(deploymentsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(deploymentsTable.AllColumns)
	query, args := stmt.Sql()

	var model Deployment
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("deployment not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating deployment", err.Error()))
	}

	deployment = model.ToEntity()
	if deployment == nil {
		return nil, errs.NewInternalServerError("failed to convert deployment model to entity", nil)
	}

	return deployment, nil
}
