package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ProjectUsecase interface {
	CreateProject(ctx context.Context, project CreateProjectInput) (*entity.Project, error)
	FindOneProject(ctx context.Context, id FindOneProjectInput) (*entity.Project, error)
	FindAllProjects(ctx context.Context, input FindAllProjectsInput) (entity.Page[entity.Project], error)
	UpdateProject(ctx context.Context, input UpdateProjectInput) (*entity.Project, error)
	DeleteProject(ctx context.Context, id DeleteProjectInput) (*entity.Project, error)
}

type projectUsecase struct {
	appConfig         config.AppConfig
	projectRepository repository.ProjectRepository
	userRepository    repository.UserRepository
}

func NewProjectUsecase(
	appConfig config.AppConfig,
	projectRepository repository.ProjectRepository,
	userRepository repository.UserRepository,
) ProjectUsecase {
	return &projectUsecase{
		appConfig:         appConfig,
		projectRepository: projectRepository,
		userRepository:    userRepository,
	}
}
