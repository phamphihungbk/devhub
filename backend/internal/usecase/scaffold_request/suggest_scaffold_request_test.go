package usecase

import (
	"context"
	"testing"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type fakePluginRepository struct {
	plugins entity.Plugins
}

func (r fakePluginRepository) CreateOne(ctx context.Context, plugin *entity.Plugin) (*entity.Plugin, error) {
	return plugin, nil
}

func (r fakePluginRepository) FindOne(ctx context.Context, id uuid.UUID) (*entity.Plugin, error) {
	return nil, nil
}

func (r fakePluginRepository) FindAll(ctx context.Context, filter repository.FindAllPluginsFilter) (*entity.Plugins, int64, error) {
	return &r.plugins, int64(len(r.plugins)), nil
}

func (r fakePluginRepository) UpdateOne(ctx context.Context, input repository.UpdatePluginInput) (*entity.Plugin, error) {
	return nil, nil
}

func (r fakePluginRepository) DeleteOne(ctx context.Context, id uuid.UUID) (*entity.Plugin, error) {
	return nil, nil
}

func TestSuggestScaffoldRequestSelectsPluginAndVariablesFromPrompt(t *testing.T) {
	scaffolderID := uuid.New()
	usecase := &scaffoldRequestUsecase{
		pluginRepository: fakePluginRepository{
			plugins: entity.Plugins{
				{
					ID:          uuid.New(),
					Name:        "Node Worker",
					Type:        entity.PluginScaffolder,
					Runtime:     entity.PluginRuntimeNode,
					Enabled:     true,
					Description: "Node worker scaffold",
					InstalledAt: time.Now(),
				},
				{
					ID:          scaffolderID,
					Name:        "Go HTTP API",
					Type:        entity.PluginScaffolder,
					Runtime:     entity.PluginRuntimeGo,
					Enabled:     true,
					Description: "Go HTTP API scaffold with Postgres support",
					InstalledAt: time.Now(),
				},
			},
		},
	}

	suggestion, err := usecase.SuggestScaffoldRequest(context.Background(), SuggestScaffoldRequestInput{
		ProjectID:          uuid.NewString(),
		Prompt:             "Create a Go payment API with Postgres and structured logging",
		ProjectName:        "Platform",
		ProjectDescription: "Handles payment workflows",
		Environments:       []string{"dev", "prod"},
	})
	if err != nil {
		t.Fatalf("SuggestScaffoldRequest returned error: %v", err)
	}

	if suggestion.PluginID != scaffolderID.String() {
		t.Fatalf("expected Go plugin, got %q", suggestion.PluginID)
	}
	if suggestion.Environment != "dev" {
		t.Fatalf("expected first project environment, got %q", suggestion.Environment)
	}
	if suggestion.Variables.ServiceName != "payment" {
		t.Fatalf("expected service name from prompt, got %q", suggestion.Variables.ServiceName)
	}
	if suggestion.Variables.Database != "postgres" {
		t.Fatalf("expected postgres database, got %q", suggestion.Variables.Database)
	}
	if !suggestion.Variables.EnableLogging {
		t.Fatal("expected logging to be enabled")
	}
}

func TestSuggestScaffoldRequestInfersDatabaseAndLoggingFromPrompt(t *testing.T) {
	usecase := &scaffoldRequestUsecase{
		pluginRepository: fakePluginRepository{
			plugins: entity.Plugins{
				{
					ID:          uuid.New(),
					Name:        "Generic Scaffolder",
					Type:        entity.PluginScaffolder,
					Runtime:     entity.PluginRuntimeGo,
					Enabled:     true,
					Description: "Generic scaffold",
					InstalledAt: time.Now(),
				},
			},
		},
	}

	suggestion, err := usecase.SuggestScaffoldRequest(context.Background(), SuggestScaffoldRequestInput{
		ProjectID: uuid.NewString(),
		Prompt:    "Create customer document search backed by MongoDB with no logging",
	})
	if err != nil {
		t.Fatalf("SuggestScaffoldRequest returned error: %v", err)
	}

	if suggestion.Variables.Database != "mongodb" {
		t.Fatalf("expected mongodb database, got %q", suggestion.Variables.Database)
	}
	if suggestion.Variables.EnableLogging {
		t.Fatal("expected logging to be disabled")
	}
}

func TestSuggestScaffoldRequestInfersExplicitPromptValues(t *testing.T) {
	usecase := &scaffoldRequestUsecase{
		pluginRepository: fakePluginRepository{
			plugins: entity.Plugins{
				{
					ID:          uuid.New(),
					Name:        "Go HTTP API",
					Type:        entity.PluginScaffolder,
					Runtime:     entity.PluginRuntimeGo,
					Enabled:     true,
					Description: "Go HTTP API scaffold",
					InstalledAt: time.Now(),
				},
			},
		},
	}

	suggestion, err := usecase.SuggestScaffoldRequest(context.Background(), SuggestScaffoldRequestInput{
		ProjectID: uuid.NewString(),
		Prompt:    "create Go payment with mysql database and port 8000 with dev, prod and staging environment",
	})
	if err != nil {
		t.Fatalf("SuggestScaffoldRequest returned error: %v", err)
	}

	if suggestion.Variables.ServiceName != "payment" {
		t.Fatalf("expected payment service name, got %q", suggestion.Variables.ServiceName)
	}
	if suggestion.Variables.Database != "mysql" {
		t.Fatalf("expected mysql database, got %q", suggestion.Variables.Database)
	}
	if suggestion.Variables.Port != 8000 {
		t.Fatalf("expected explicit port 8000, got %d", suggestion.Variables.Port)
	}
	if suggestion.Environment != "dev" {
		t.Fatalf("expected first mentioned environment dev, got %q", suggestion.Environment)
	}
	expectedEnvironments := []string{"dev", "staging", "prod"}
	if len(suggestion.Environments) != len(expectedEnvironments) {
		t.Fatalf("expected environments %v, got %v", expectedEnvironments, suggestion.Environments)
	}
	for index, expected := range expectedEnvironments {
		if suggestion.Environments[index] != expected {
			t.Fatalf("expected environments %v, got %v", expectedEnvironments, suggestion.Environments)
		}
	}
}
