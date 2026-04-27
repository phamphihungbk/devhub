package usecase

import (
	"context"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/ai"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type SuggestScaffoldRequestInput struct {
	ProjectID string `json:"project_id" validate:"required,uuid"`
	Prompt    string `json:"prompt" validate:"required"`
}

type ScaffoldRequestSuggestion struct {
	Source      string                          `json:"source"`
	PluginID    string                          `json:"plugin_id"`
	PluginName  string                          `json:"plugin_name"`
	Confidence  float64                         `json:"confidence"`
	Environment string                          `json:"environment"`
	Variables   entity.ScaffoldRequestVariables `json:"variables"`
	Rationale   []string                        `json:"rationale"`
}

func (u *scaffoldRequestUsecase) SuggestScaffoldRequest(ctx context.Context, input SuggestScaffoldRequestInput) (suggestion ScaffoldRequestSuggestion, err error) {
	const errLocation = "[usecase scaffold_request/suggest_scaffold_request SuggestScaffoldRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)

	if err != nil {
		return ScaffoldRequestSuggestion{}, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return ScaffoldRequestSuggestion{}, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	projectID := uuid.MustParse(input.ProjectID)
	project, err := u.projectRepository.FindOne(ctx, projectID)
	if err != nil {
		return ScaffoldRequestSuggestion{}, misc.WrapError(err, errs.NewInternalServerError("failed to load project context", nil))
	}

	projectEnvironments := make([]string, 0, len(project.Environments))
	for _, environment := range project.Environments {
		projectEnvironments = append(projectEnvironments, environment.String())
	}

	plugin, plan, err := u.suggestScaffolderPlugin(ctx, input.Prompt+" "+project.Description)
	if err != nil {
		return ScaffoldRequestSuggestion{}, err
	}

	aiSuggestion, err := u.aiSuggestionGenerator.GenerateScaffoldSuggestion(ctx, ai.ScaffoldSuggestionInput{
		Prompt:              input.Prompt,
		Project:             *project,
		ProjectEnvironments: projectEnvironments,
		Plugin:              *plugin,
		Plan:                *plan,
	})
	if err != nil {
		return ScaffoldRequestSuggestion{}, misc.WrapError(err, errs.NewInternalServerError("failed to generate scaffold suggestion", nil))
	}

	return ScaffoldRequestSuggestion{
		Source:      aiSuggestion.Source,
		PluginID:    plugin.ID.String(),
		PluginName:  plugin.Name,
		Confidence:  plan.Confidence,
		Environment: aiSuggestion.Environment,
		Variables:   aiSuggestion.Variables,
		Rationale:   aiSuggestion.Rationale,
	}, nil
}

func (u *scaffoldRequestUsecase) suggestScaffolderPlugin(ctx context.Context, intent string) (*entity.Plugin, *ai.ScaffoldPlan, error) {
	plugins, _, err := u.pluginRepository.FindAll(ctx, repository.FindAllPluginsFilter{})
	if err != nil {
		return nil, nil, misc.WrapError(err, errs.NewInternalServerError("failed to load scaffold plugins", nil))
	}
	if plugins == nil {
		return nil, nil, errs.NewBadRequestError("no scaffold plugins are available", nil)
	}

	enabledPlugins := make([]entity.Plugin, 0, len(*plugins))
	candidates := make([]ai.PluginCandidate, 0, len(*plugins))
	for index := range *plugins {
		plugin := (*plugins)[index]
		if !plugin.Enabled || plugin.Type != entity.PluginScaffolder {
			continue
		}
		enabledPlugins = append(enabledPlugins, plugin)
		candidates = append(candidates, newScaffoldPluginCandidate(plugin))
	}

	if len(enabledPlugins) == 0 {
		return nil, nil, errs.NewBadRequestError("no enabled scaffolder plugins are available", nil)
	}

	plan, err := u.aiClient.PlanScaffold(ctx, ai.ScaffoldPlanningInput{
		Prompt:  intent,
		Plugins: candidates,
	})

	if err != nil {
		return nil, nil, misc.WrapError(err, errs.NewInternalServerError("failed to rank scaffold plugins", nil))
	}

	for index := range enabledPlugins {
		if enabledPlugins[index].Name == plan.PluginName {
			return &enabledPlugins[index], plan, nil
		}
	}

	return nil, nil, errs.NewBadRequestError("ranked scaffold plugin is not available", nil)
}

func newScaffoldPluginCandidate(plugin entity.Plugin) ai.PluginCandidate {
	return ai.PluginCandidate{
		ID:           plugin.ID.String(),
		Name:         plugin.Name,
		Type:         plugin.Type.String(),
		Runtime:      plugin.Runtime.String(),
		Entrypoint:   plugin.Entrypoint,
		Description:  plugin.Description,
		Keywords:     inferScaffoldPluginKeywords(plugin),
		Capabilities: inferScaffoldPluginCapabilities(plugin),
	}
}

func inferScaffoldPluginKeywords(plugin entity.Plugin) []string {
	value := strings.ToLower(plugin.Name + " " + plugin.Description + " " + plugin.Entrypoint)
	keywords := []string{}
	for _, keyword := range []string{
		"api", "backend", "database", "fastapi", "frontend", "go", "golang", "grpc", "http", "job", "mysql", "node", "postgres", "python", "react", "rest", "vue", "worker",
	} {
		if strings.Contains(value, keyword) {
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}

func inferScaffoldPluginCapabilities(plugin entity.Plugin) []string {
	value := strings.ToLower(plugin.Name + " " + plugin.Description + " " + plugin.Entrypoint)
	capabilities := make([]string, 0)
	for _, capability := range []string{
		"api service", "background worker", "frontend application", "grpc service", "http service", "web application",
	} {
		for _, token := range strings.Fields(capability) {
			if strings.Contains(value, token) {
				capabilities = append(capabilities, capability)
				break
			}
		}
	}
	return capabilities
}
