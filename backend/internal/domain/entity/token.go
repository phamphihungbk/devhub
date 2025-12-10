package entity

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken struct {
	UserID    uuid.UUID
	Role      UserRole
	Issuer    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type AccessTokenClaims struct {
	Role string
	jwt.RegisteredClaims
}

func (t AccessToken) ToJWTClaims() AccessTokenClaims {
	return AccessTokenClaims{
		Role: t.Role.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.Issuer,
			Subject:   t.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(t.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(t.IssuedAt),
		},
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
