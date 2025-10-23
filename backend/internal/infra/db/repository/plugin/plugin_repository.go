package repository

import (
	"context"

	"github.com/phamphihungbk/devhub-backend/internal/app/models"
	"gorm.io/gorm"
)

type PluginRepository struct {
	db *gorm.DB
}

type PluginRepositoryInterface interface {
	Create(ctx context.Context, plugin *models.Plugin) error
	GetByID(ctx context.Context, id string) (*models.Plugin, error)
	List(ctx context.Context) ([]models.Plugin, error)
	Update(ctx context.Context, plugin *models.Plugin) error
	Delete(ctx context.Context, id string) error
}

func NewPluginRepository(db *gorm.DB) *PluginRepository {
	return &PluginRepository{db: db}
}

func (r *PluginRepository) Create(ctx context.Context, plugin *models.Plugin) error {
	return r.db.WithContext(ctx).Create(plugin).Error
}

func (r *PluginRepository) GetByID(ctx context.Context, id string) (*models.Plugin, error) {
	var plugin models.Plugin
	err := r.db.WithContext(ctx).First(&plugin, "id = ?", id).Error
	return &plugin, err
}

func (r *PluginRepository) List(ctx context.Context) ([]models.Plugin, error) {
	var plugins []models.Plugin
	err := r.db.WithContext(ctx).Find(&plugins).Error
	return plugins, err
}

func (r *PluginRepository) Update(ctx context.Context, plugin *models.Plugin) error {
	return r.db.WithContext(ctx).Save(plugin).Error
}

func (r *PluginRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Plugin{}, "id = ?", id).Error
}
