package repository

import (
	"context"

	"github.com/phamphihungbk/devhub-backend/internal/app/models"
	"gorm.io/gorm"
)

type DeploymentRepository struct {
	db *gorm.DB
}

type DeploymentRepositoryInterface interface {
	Create(ctx context.Context, deployment *models.Deployment) error
	GetByID(ctx context.Context, id string) (*models.Deployment, error)
	List(ctx context.Context) ([]models.Deployment, error)
	Update(ctx context.Context, deployment *models.Deployment) error
	Delete(ctx context.Context, id string) error
}

func NewDeploymentRepository(db *gorm.DB) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

func (r *DeploymentRepository) Create(ctx context.Context, deployment *models.Deployment) error {
	return r.db.WithContext(ctx).Create(deployment).Error
}

func (r *DeploymentRepository) GetByID(ctx context.Context, id string) (*models.Deployment, error) {
	var deployment models.Deployment
	err := r.db.WithContext(ctx).First(&deployment, "id = ?", id).Error
	return &deployment, err
}

func (r *DeploymentRepository) List(ctx context.Context) ([]models.Deployment, error) {
	var deployments []models.Deployment
	err := r.db.WithContext(ctx).Find(&deployments).Error
	return deployments, err
}

func (r *DeploymentRepository) Update(ctx context.Context, deployment *models.Deployment) error {
	return r.db.WithContext(ctx).Save(deployment).Error
}

func (r *DeploymentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Deployment{}, "id = ?", id).Error
}
