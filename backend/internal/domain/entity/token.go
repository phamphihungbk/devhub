package entity

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken struct {
	UserID    uuid.UUID
	Role      UserRole
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (t AccessToken) ToJWTClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub":  t.UserID.String(),
		"role": t.Role.String(),
		"exp":  t.ExpiresAt.Unix(),
		"iat":  t.IssuedAt.Unix(),
		"iss":  "devhub",
	}
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	DeletedAt time.Time
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
