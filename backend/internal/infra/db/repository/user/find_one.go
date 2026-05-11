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

func (r *userRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (user *entity.User, err error) {
	const errLocation = "[repository user/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	usersTable := table.Users
	userRolesTable := table.UserRoles
	rolesTable := table.Roles

	// SQL statement
	stmt := postgres.SELECT(
		usersTable.AllColumns,
		postgres.RawString("COALESCE(jsonb_agg(roles.name ORDER BY roles.name) FILTER (WHERE roles.id IS NOT NULL), '[]'::jsonb)").AS("roles_json"),
	).
		FROM(
			usersTable.
				LEFT_JOIN(userRolesTable, userRolesTable.UserID.EQ(usersTable.ID)).
				LEFT_JOIN(rolesTable, rolesTable.ID.EQ(userRolesTable.RoleID)),
		).
		WHERE(usersTable.ID.EQ(postgres.UUID(id))).
		GROUP_BY(usersTable.ID)

	query, args := stmt.Sql()

	var model User
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("user not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying user by id", err.Error()))
	}

	user = model.ToEntity()
	if user == nil {
		return nil, errs.NewInternalServerError("failed to convert user model to entity", nil)
	}

	return user, nil
}
