package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Role      UserRole       `gorm:"type:varchar(16);default:'user'" json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	LastLogin time.Time      `json:"lastLogin"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
