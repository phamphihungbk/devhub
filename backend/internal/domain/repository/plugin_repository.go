package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type PluginRepository interface {
	CreateOne(ctx context.Context, plugin *entity.Plugin) (*entity.Plugin, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Plugin, error)
	FindAll(ctx context.Context, filter FindAllPluginsFilter) (*entity.Plugins, int64, error)
	UpdateOne(ctx context.Context, input UpdatePluginInput) (*entity.Plugin, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.Plugin, error)
}

type FindAllPluginsFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}

type UpdatePluginInput struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Type        *entity.PluginType
	Version     *string
	Runtime     *entity.PluginRuntime
	Entrypoint  *string
	Scope       *entity.PluginScope
	Enabled     *bool
}
