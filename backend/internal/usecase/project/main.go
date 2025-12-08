package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user CreateUserInput) (*entity.User, error)
	FindOneUser(ctx context.Context, id FindOneUserInput) (*entity.User, error)
	FindAllUsers(ctx context.Context, input FindAllUsersInput) (entity.Page[entity.User], error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (*entity.User, error)
	DeleteUser(ctx context.Context, id DeleteUserInput) (*entity.User, error)
}

type userUsecase struct {
	appConfig      config.AppConfig
	userRepository repository.UserRepository
}

func NewUserUsecase(
	appConfig config.AppConfig,
	userRepository repository.UserRepository,
) UserUsecase {
	return &userUsecase{
		appConfig:      appConfig,
		userRepository: userRepository,
	}
}
