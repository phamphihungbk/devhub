package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserRole = fmt.Errorf("invalid user role")
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

var userRoleStringMapper = map[UserRole]string{
	RoleAdmin: "admin",
	RoleUser:  "user",
}

func (s UserRole) String() string {
	return userRoleStringMapper[s]
}

func (s UserRole) IsValid() bool {
	switch s {
	case RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}

// Parse parses a string into a UserRole. It returns an error if the string is not a valid UserRole.
func (s UserRole) Parse(role string) (UserRole, error) {
	userRole := UserRole(role)
	if !userRole.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidUserRole, role)
	}
	return userRole, nil
}

func (s UserRole) MustParse(role string) UserRole {
	userRole := UserRole(role)
	if !userRole.IsValid() {
		panic(`user role: Parse(` + s + `): `)
	}
	return userRole
}

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type Users []User
