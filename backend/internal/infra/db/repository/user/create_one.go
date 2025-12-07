package userrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *userRepositoryImpl) CreateOne(ctx context.Context, input *entity.User) (user *entity.User, err error) {
	const errLocation = "[repository user/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	usersTable := table.Users
	// SQL statement
	stmt := usersTable.INSERT(
		usersTable.AllColumns.Except(usersTable.DefaultColumns), // Exclude columns with default values
	).MODEL(model.Users{
		Name:  input.Name,
		Email: input.Email,
		Role:  string(input.Role),
	}).RETURNING(usersTable.AllColumns)
	query, args := stmt.Sql()

	var model User
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating concert", err.Error()))
	}

	user = model.ToEntity()
	if user == nil {
		return nil, errs.NewInternalServerError("failed to convert user model to entity", nil)
	}

	return
}
