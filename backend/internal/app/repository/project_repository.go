package repository

import (
	"context"

	"github.com/phamphihungbk/devhub-backend/internal/app/models"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

type ProjectRepositoryInterface interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id string) (*models.Project, error)
	List(ctx context.Context) ([]models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id string) error
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error
	return &project, err
}

func (r *ProjectRepository) List(ctx context.Context) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.WithContext(ctx).Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id).Error
}
