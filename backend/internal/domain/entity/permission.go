package entity

import (
	"time"

	"github.com/google/uuid"
)

type PermissionName string

const (
	PermissionUserRead             PermissionName = "user.read"
	PermissionUserWrite            PermissionName = "user.write"
	PermissionProjectWrite         PermissionName = "project.write"
	PermissionScaffoldRequestWrite PermissionName = "scaffold_request.write"
	PermissionReleaseWrite         PermissionName = "release.write"
	PermissionDeploymentWrite      PermissionName = "deployment.write"
	PermissionPluginWrite          PermissionName = "plugin.write"
)

func (p PermissionName) String() string {
	return string(p)
}

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}

type Permissions []Permission
