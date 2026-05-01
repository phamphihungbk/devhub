package userrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *userRepositoryImpl) HasPermissions(ctx context.Context, userID uuid.UUID, permissions []entity.PermissionName) (bool, error) {
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
		permissionsTable.Name,
	).
		FROM(
			userRolesTable.
				INNER_JOIN(rolePermissionsTable, rolePermissionsTable.RoleID.EQ(userRolesTable.RoleID)).
				INNER_JOIN(permissionsTable, permissionsTable.ID.EQ(rolePermissionsTable.PermissionID)),
		).
		WHERE(
			userRolesTable.UserID.EQ(postgres.UUID(userID)).
				AND(permissionsTable.Name.IN(permissionExpressions...)),
		).
		GROUP_BY(permissionsTable.Name)

	query, args := stmt.Sql()

	var matchedPermissions []struct {
		Name string `db:"name"`
	}
	if err := r.execer.SelectContext(ctx, &matchedPermissions, query, args...); err != nil {
		return false, err
	}

	return len(matchedPermissions) == len(requiredPermissions), nil
}
