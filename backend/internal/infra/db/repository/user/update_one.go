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

func (r *userRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateUserInput) (user *entity.User, err error) {
	const errLocation = "[repository user/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	usersTable := table.Users
	var updateModel User
	columns := make(postgres.ColumnList, 0)

	// build the update model
	if input.Name != nil {
		updateModel.Name = misc.GetValue(input.Name)
		columns = append(columns, usersTable.Name)
	}

	if input.Role != nil {
		updateModel.Role = string(misc.GetValue(input.Role))
		columns = append(columns, usersTable.Role)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	// SQL statement
	stmt := usersTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(usersTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(usersTable.AllColumns)
	query, args := stmt.Sql()

	var model User
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("user not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating user", err.Error()))
	}

	user = model.ToEntity()
	if user == nil {
		return nil, errs.NewInternalServerError("failed to convert user model to entity", nil)
	}

	return user, nil
}
