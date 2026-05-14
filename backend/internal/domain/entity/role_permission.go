package entity

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	ID           uuid.UUID
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	CreatedAt    time.Time
}

type RolePermissions []RolePermission
