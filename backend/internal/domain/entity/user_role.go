package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserAssignedRole string

const (
	RolePlatformAdmin UserAssignedRole = "platform_admin"
	RoleOrgAdmin      UserAssignedRole = "org_admin"
	RoleTeamLead      UserAssignedRole = "team_lead"
	RoleDeveloper     UserAssignedRole = "developer"
	RoleViewer        UserAssignedRole = "viewer"
)

func (u UserAssignedRole) String() string {
	return string(u)
}

type UserRole struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	RoleID    uuid.UUID
	CreatedAt time.Time
}

type UserRoles []UserRole
