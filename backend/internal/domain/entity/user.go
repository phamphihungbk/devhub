package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Role      UserRole
	CreatedAt time.Time
	LastLogin time.Time
	DeletedAt time.Time
}

type Users []User
