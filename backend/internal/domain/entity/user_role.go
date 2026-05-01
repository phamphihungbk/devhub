package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	RoleID    uuid.UUID
	CreatedAt time.Time
}

type UserRoles []UserRole
