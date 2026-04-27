package ai

import "context"

type Client interface {
	PlanScaffold(ctx context.Context, input ScaffoldPlanningInput) (*ScaffoldPlan, error)
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
