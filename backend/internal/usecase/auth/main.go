package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type AuthUsecase interface {
	LoginUser(ctx context.Context, input LoginUserInput) (*entity.User, error)
	LogoutUser(ctx context.Context) (*entity.User, error)
}

type authUsecase struct {
	appConfig      config.AppConfig
	userRepository repository.UserRepository
}

func NewAuthUsecase(
	appConfig config.AppConfig,
	userRepository repository.UserRepository,
) AuthUsecase {
	return &authUsecase{
		appConfig:      appConfig,
		userRepository: userRepository,
	}
}
