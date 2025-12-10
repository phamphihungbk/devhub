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

func (r *userRepositoryImpl) DeleteOne(ctx context.Context, id uuid.UUID) (user *entity.User, err error) {
	const errLocation = "[repository user/delete Delete] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	usersTable := table.Users

	// SQL statement
	stmt := usersTable.DELETE().
		WHERE(usersTable.ID.EQ(postgres.UUID(id))).
		RETURNING(usersTable.AllColumns)
	query, args := stmt.Sql()

	var model User
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("user not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while deleting user", err.Error()))
	}

	user = model.ToEntity()
	if user == nil {
		return nil, errs.NewInternalServerError("failed to convert user model to entity", nil)
	}

	return user, nil

}
