package userrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *userRepositoryImpl) HasPermissions(ctx context.Context, userID uuid.UUID, permissions []entity.PermissionName) (allowed bool, err error) {
	const errLocation = "[repository user/has_permissions HasPermissions] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	if len(permissions) == 0 {
		return true, nil
	}

	requiredPermissions := make([]string, 0, len(permissions))
	permissionExpressions := make([]postgres.Expression, 0, len(permissions))
	for _, permission := range permissions {
		permissionName := string(permission)
		requiredPermissions = append(requiredPermissions, permissionName)
		permissionExpressions = append(permissionExpressions, postgres.String(permissionName))
	}

	userRolesTable := table.UserRoles
	rolePermissionsTable := table.RolePermissions
	permissionsTable := table.Permissions

	stmt := postgres.SELECT(
		postgres.COUNT(postgres.DISTINCT(permissionsTable.Name)).AS("count"),
	).
		FROM(
			userRolesTable.
				INNER_JOIN(rolePermissionsTable, rolePermissionsTable.RoleID.EQ(userRolesTable.RoleID)).
				INNER_JOIN(permissionsTable, permissionsTable.ID.EQ(rolePermissionsTable.PermissionID)),
		).
		WHERE(
			userRolesTable.UserID.EQ(postgres.UUID(userID)).
				AND(permissionsTable.Name.IN(permissionExpressions...)),
		)

	query, args := stmt.Sql()

	var result struct {
		Count int64 `db:"count"`
	}
	if err := r.execer.GetContext(ctx, &result, query, args...); err != nil {
		return false, misc.WrapError(err, errs.NewDatabaseError("error while checking user permissions", err.Error()))
	}

	return result.Count == int64(len(requiredPermissions)), nil
}
