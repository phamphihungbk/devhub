package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Roles        []string
	TeamID       uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type Users []User
