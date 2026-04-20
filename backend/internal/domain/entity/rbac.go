package entity

type Permission string

const (
	PermissionUserRead             Permission = "user.read"
	PermissionUserWrite            Permission = "user.write"
	PermissionProjectWrite         Permission = "project.write"
	PermissionScaffoldRequestWrite Permission = "scaffold_request.write"
	PermissionReleaseWrite         Permission = "release.write"
	PermissionDeploymentWrite      Permission = "deployment.write"
	PermissionPluginWrite          Permission = "plugin.write"
)

var rolePermissions = map[UserRole]map[Permission]struct{}{
	RolePlatformAdmin: {
		PermissionUserRead:             {},
		PermissionUserWrite:            {},
		PermissionProjectWrite:         {},
		PermissionScaffoldRequestWrite: {},
		PermissionReleaseWrite:         {},
		PermissionDeploymentWrite:      {},
		PermissionPluginWrite:          {},
	},
	RoleOrgAdmin: {
		PermissionUserRead:             {},
		PermissionUserWrite:            {},
		PermissionProjectWrite:         {},
		PermissionScaffoldRequestWrite: {},
		PermissionReleaseWrite:         {},
		PermissionDeploymentWrite:      {},
	},
	RoleTeamLead: {
		PermissionProjectWrite:         {},
		PermissionScaffoldRequestWrite: {},
		PermissionReleaseWrite:         {},
		PermissionDeploymentWrite:      {},
	},
	RoleDeveloper: {
		PermissionScaffoldRequestWrite: {},
		PermissionReleaseWrite:         {},
		PermissionDeploymentWrite:      {},
	},
	RoleViewer: {},
}

func (s UserRole) HasPermissions(permissions ...Permission) bool {
	allowedPermissions, exists := rolePermissions[s]
	if !exists {
		return false
	}

	for _, permission := range permissions {
		if _, ok := allowedPermissions[permission]; !ok {
			return false
		}
	}

	return true
}
