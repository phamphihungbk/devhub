package ai

import (
	"context"

	"devhub-backend/internal/domain/entity"
)

type Client interface {
	PlanScaffold(ctx context.Context, input ScaffoldPlanningInput) (*ScaffoldPlan, error)
}

type ScaffoldSuggestionGenerator interface {
	GenerateScaffoldSuggestion(ctx context.Context, input ScaffoldSuggestionInput) (*ScaffoldSuggestion, error)
}

type ScaffoldPlanningInput struct {
	Prompt  string            `json:"prompt"`
	Plugins []PluginCandidate `json:"plugins"`
}

type PluginCandidate struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Runtime      string         `json:"runtime"`
	Entrypoint   string         `json:"entrypoint"`
	Description  string         `json:"description"`
	Keywords     []string       `json:"keywords"`
	Capabilities []string       `json:"capabilities"`
	Schema       map[string]any `json:"schema"`
}

type ScaffoldPlan struct {
	PluginName string         `json:"plugin_name"`
	Variables  map[string]any `json:"variables"`
	Confidence float64        `json:"confidence"`
	Reason     string         `json:"reason"`
	Matches    []string       `json:"matches"`
}

type ScaffoldSuggestionInput struct {
	Prompt              string
	Project             entity.Project
	ProjectEnvironments []string
	Plugin              entity.Plugin
	Plan                ScaffoldPlan
}

type ScaffoldSuggestion struct {
	Source      string
	Environment string
	Variables   entity.ScaffoldRequestVariables
	Rationale   []string
}
