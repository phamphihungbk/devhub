package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type AuthUsecase interface {
	IssueToken(ctx context.Context, input IssueTokenInput) (*entity.Token, error)
	RevokeToken(ctx context.Context, input RevokeTokenInput) (*entity.RefreshToken, error)
	FindOneUser(ctx context.Context, input FindOneUserInput) (*entity.User, error)
}

type authUsecase struct {
	tokenConfig            config.TokenConfig
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
}

func NewAuthUsecase(
	tokenConfig config.TokenConfig,
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
) AuthUsecase {
	return &authUsecase{
		tokenConfig:            tokenConfig,
		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
	}
}
