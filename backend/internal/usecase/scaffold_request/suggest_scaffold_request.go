package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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
	Source       string                          `json:"source"`
	PluginID     string                          `json:"plugin_id"`
	PluginName   string                          `json:"plugin_name"`
	Confidence   float64                         `json:"confidence"`
	Environment  string                          `json:"environment"`
	Environments []string                        `json:"environments"`
	Variables    entity.ScaffoldRequestVariables `json:"variables"`
	Rationale    []string                        `json:"rationale"`
}

var (
	nonScaffoldServiceNameCharacterPattern = regexp.MustCompile(`[^a-z0-9]+`)
	portPattern                            = regexp.MustCompile(`\bport\s*[:=]?\s*(\d{1,5})\b`)
)

func (u *scaffoldRequestUsecase) SuggestScaffoldRequest(ctx context.Context, input SuggestScaffoldRequestInput) (suggestion *ScaffoldRequestSuggestion, err error) {
	const errLocation = "[usecase scaffold_request/suggest_scaffold_request SuggestScaffoldRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	projectID := uuid.MustParse(input.ProjectID)
	project, err := u.projectRepository.FindOne(ctx, projectID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to load project context", nil))
	}

	projectEnvironments := make([]string, 0, len(project.Environments))
	for _, environment := range project.Environments {
		projectEnvironments = append(projectEnvironments, environment.String())
	}

	plugin, plan, err := u.suggestScaffolderPlugin(ctx, input.Prompt+" "+project.Description)
	if err != nil {
		return nil, err
	}

	serviceName := normalizeScaffoldSuggestionServiceName(firstNonEmpty(inferNameFromPrompt(input.Prompt), project.Name, inferNameFromPrompt(project.Description), "new-service"))
	suggestedEnvironments := inferScaffoldSuggestionEnvironments(input.Prompt, projectEnvironments)
	environment := pickScaffoldSuggestionEnvironment("", suggestedEnvironments)
	modulePath := suggestScaffoldSuggestionModulePath(plugin, project.Name, serviceName)
	database := suggestScaffoldSuggestionDatabase(input.Prompt, project.Description, plugin)
	port := suggestScaffoldSuggestionPort(input.Prompt, serviceName)
	enableLogging := suggestScaffoldSuggestionLogging(input.Prompt, project.Description)

	rationale := []string{
		"Prompt was analyzed by the local token ranking client.",
		fmt.Sprintf("Winning plugin score: %.0f%%.", plan.Confidence*100),
		"Scaffolder plugin was selected from enabled scaffold_request plugins only.",
		"Service name was inferred from the user prompt and project context.",
		"Selected plugin: " + plugin.Name + ".",
	}
	if len(plan.Matches) > 0 {
		rationale = append(rationale, "Matched tokens: "+strings.Join(plan.Matches, ", ")+".")
	}
	if hasExplicitPort(input.Prompt) {
		rationale = append(rationale, "Port was taken from the prompt.")
	} else {
		rationale = append(rationale, "Port was selected from a stable service-name hash.")
	}
	if len(suggestedEnvironments) > 0 {
		rationale = append(rationale, "Environments were inferred from the prompt or project settings.")
	}
	if strings.TrimSpace(project.Description) != "" {
		rationale = append(rationale, "Project description was used as additional intent context.")
	}

	return ScaffoldRequestSuggestion{
		Source:       "local-token-ranker-v1",
		PluginID:     plugin.ID.String(),
		PluginName:   plugin.Name,
		Confidence:   plan.Confidence,
		Environment:  environment,
		Environments: suggestedEnvironments,
		Variables: entity.ScaffoldRequestVariables{
			ServiceName:   serviceName,
			ModulePath:    modulePath,
			Port:          port,
			Database:      database,
			EnableLogging: enableLogging,
		},
		Rationale: rationale,
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func inferNameFromPrompt(prompt string) string {
	words := strings.Fields(strings.ToLower(prompt))
	if len(words) == 0 {
		return ""
	}

	stopWords := map[string]struct{}{
		"a": {}, "an": {}, "and": {}, "api": {}, "app": {}, "build": {}, "create": {}, "for": {}, "generate": {},
		"database": {}, "db": {}, "dev": {}, "environment": {}, "go": {}, "http": {}, "logging": {}, "mariadb": {}, "mongodb": {},
		"mysql": {}, "node": {}, "of": {}, "please": {}, "port": {}, "postgres": {}, "prod": {}, "python": {},
		"redis": {}, "scaffold": {}, "service": {}, "staging": {}, "structured": {}, "the": {}, "to": {}, "with": {},
	}

	selected := make([]string, 0, 3)
	for _, word := range words {
		word = strings.Trim(word, ".,:;()[]{}")
		if _, ok := stopWords[word]; ok || word == "" {
			continue
		}
		if _, err := strconv.Atoi(word); err == nil {
			continue
		}
		selected = append(selected, word)
		if len(selected) == 3 {
			break
		}
	}

	return strings.Join(selected, "-")
}

func normalizeScaffoldSuggestionServiceName(value string) string {
	normalized := strings.Trim(nonScaffoldServiceNameCharacterPattern.ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "-"), "-")
	if normalized == "" {
		return "new-service"
	}
	return normalized
}

func pickScaffoldSuggestionEnvironment(selected string, environments []string) string {
	selected = strings.TrimSpace(selected)
	if selected != "" {
		return selected
	}
	for _, environment := range environments {
		if strings.TrimSpace(environment) != "" {
			return environment
		}
	}
	return entity.EnvDev.String()
}

func inferScaffoldSuggestionEnvironments(prompt string, projectEnvironments []string) []string {
	value := strings.ToLower(prompt)
	known := []string{entity.EnvDev.String(), entity.EnvStaging.String(), entity.EnvProd.String()}
	seen := map[string]struct{}{}
	environments := make([]string, 0, len(known))

	for _, environment := range known {
		if strings.Contains(value, environment) {
			seen[environment] = struct{}{}
			environments = append(environments, environment)
		}
	}

	if len(environments) > 0 {
		return environments
	}

	for _, environment := range projectEnvironments {
		environment = strings.TrimSpace(environment)
		if environment == "" {
			continue
		}
		if _, ok := seen[environment]; ok {
			continue
		}
		seen[environment] = struct{}{}
		environments = append(environments, environment)
	}

	if len(environments) == 0 {
		return []string{entity.EnvDev.String()}
	}

	return environments
}

func suggestScaffoldSuggestionModulePath(plugin *entity.Plugin, projectName string, serviceName string) string {
	host := "gitea.devhub.local"
	if plugin != nil && strings.Contains(plugin.Entrypoint, "github") {
		host = "github.com"
	}
	owner := normalizeScaffoldSuggestionServiceName(firstNonEmpty(projectName, "platform"))
	return host + "/" + owner + "/" + serviceName
}

func suggestScaffoldSuggestionPort(prompt string, serviceName string) int {
	if port, ok := inferScaffoldSuggestionPort(prompt); ok {
		return port
	}

	hash := 0
	for _, char := range serviceName {
		hash += int(char)
	}
	return 8000 + hash%1000
}

func inferScaffoldSuggestionPort(prompt string) (int, bool) {
	matches := portPattern.FindStringSubmatch(strings.ToLower(prompt))
	if len(matches) != 2 {
		return 0, false
	}

	port, err := strconv.Atoi(matches[1])
	if err != nil || port < 1 || port > 65535 {
		return 0, false
	}

	return port, true
}

func hasExplicitPort(prompt string) bool {
	_, ok := inferScaffoldSuggestionPort(prompt)
	return ok
}

func suggestScaffoldSuggestionDatabase(prompt string, projectDescription string, plugin *entity.Plugin) string {
	value := strings.ToLower(prompt + " " + projectDescription)
	if plugin != nil {
		value += " " + strings.ToLower(plugin.Description)
	}
	switch {
	case strings.Contains(value, "cache"), strings.Contains(value, "redis"):
		return "redis"
	case strings.Contains(value, "mongo"), strings.Contains(value, "document"):
		return "mongodb"
	case strings.Contains(value, "mysql"), strings.Contains(value, "mariadb"):
		return "mysql"
	default:
		return "postgres"
	}
}

func suggestScaffoldSuggestionLogging(prompt string, projectDescription string) bool {
	value := strings.ToLower(prompt + " " + projectDescription)
	return !strings.Contains(value, "disable logging") && !strings.Contains(value, "no logging")
}
