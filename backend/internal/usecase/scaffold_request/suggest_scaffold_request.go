package usecase

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
)

type SuggestScaffoldRequestInput struct {
	ProjectID          string   `json:"project_id" validate:"required,uuid"`
	Prompt             string   `json:"prompt" validate:"required"`
	ProjectName        string   `json:"project_name"`
	ProjectDescription string   `json:"project_description"`
	Environment        string   `json:"environment"`
	Environments       []string `json:"environments"`
}

type ScaffoldRequestSuggestion struct {
	Source       string                          `json:"source"`
	PluginID     string                          `json:"plugin_id"`
	PluginName   string                          `json:"plugin_name"`
	Environment  string                          `json:"environment"`
	Environments []string                        `json:"environments"`
	Variables    entity.ScaffoldRequestVariables `json:"variables"`
	Rationale    []string                        `json:"rationale"`
}

var (
	nonScaffoldServiceNameCharacterPattern = regexp.MustCompile(`[^a-z0-9]+`)
	portPattern                            = regexp.MustCompile(`\bport\s*[:=]?\s*(\d{1,5})\b`)
)

func (u *scaffoldRequestUsecase) SuggestScaffoldRequest(ctx context.Context, input SuggestScaffoldRequestInput) (suggestion ScaffoldRequestSuggestion, err error) {
	const errLocation = "[usecase scaffold_request/suggest_scaffold_request SuggestScaffoldRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	if strings.TrimSpace(input.Prompt) == "" {
		return ScaffoldRequestSuggestion{}, errs.NewBadRequestError("prompt is required", nil)
	}

	plugin, err := u.suggestScaffolderPlugin(ctx, input.Prompt+" "+input.ProjectDescription)
	if err != nil {
		return ScaffoldRequestSuggestion{}, err
	}

	serviceName := normalizeScaffoldSuggestionServiceName(firstNonEmpty(inferNameFromPrompt(input.Prompt), input.ProjectName, inferNameFromPrompt(input.ProjectDescription), "new-service"))
	suggestedEnvironments := inferScaffoldSuggestionEnvironments(input.Prompt, input.Environments)
	environment := pickScaffoldSuggestionEnvironment(input.Environment, suggestedEnvironments)
	modulePath := suggestScaffoldSuggestionModulePath(plugin, input.ProjectName, serviceName)
	database := suggestScaffoldSuggestionDatabase(input.Prompt, input.ProjectDescription, plugin)
	port := suggestScaffoldSuggestionPort(input.Prompt, serviceName)
	enableLogging := suggestScaffoldSuggestionLogging(input.Prompt, input.ProjectDescription)

	rationale := []string{
		"Prompt was analyzed by the local scaffold suggestion engine.",
		"Scaffolder plugin was selected from enabled scaffold_request plugins.",
		"Service name was inferred from the user prompt and project context.",
		"Selected plugin: " + plugin.Name + ".",
	}
	if hasExplicitPort(input.Prompt) {
		rationale = append(rationale, "Port was taken from the prompt.")
	} else {
		rationale = append(rationale, "Port was selected from a stable service-name hash.")
	}
	if len(suggestedEnvironments) > 0 {
		rationale = append(rationale, "Environments were inferred from the prompt or project settings.")
	}
	if strings.TrimSpace(input.ProjectDescription) != "" {
		rationale = append(rationale, "Project description was used as additional intent context.")
	}

	return ScaffoldRequestSuggestion{
		Source:       "local-prompt-heuristic-v2",
		PluginID:     plugin.ID.String(),
		PluginName:   plugin.Name,
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

func (u *scaffoldRequestUsecase) suggestScaffolderPlugin(ctx context.Context, intent string) (*entity.Plugin, error) {
	plugins, _, err := u.pluginRepository.FindAll(ctx, repository.FindAllPluginsFilter{})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to load scaffold plugins", nil))
	}
	if plugins == nil {
		return nil, errs.NewBadRequestError("no scaffold plugins are available", nil)
	}

	var fallback *entity.Plugin
	var best *entity.Plugin
	bestScore := -1
	intent = strings.ToLower(intent)

	for index := range *plugins {
		plugin := &(*plugins)[index]
		if !plugin.Enabled || plugin.Type != entity.PluginScaffolder {
			continue
		}
		if fallback == nil {
			fallback = plugin
		}

		score := scoreScaffolderPlugin(intent, plugin)
		if score > bestScore {
			best = plugin
			bestScore = score
		}
	}

	if best != nil {
		return best, nil
	}
	if fallback != nil {
		return fallback, nil
	}

	return nil, errs.NewBadRequestError("no enabled scaffolder plugins are available", nil)
}

func scoreScaffolderPlugin(intent string, plugin *entity.Plugin) int {
	value := strings.ToLower(plugin.Name + " " + plugin.Description + " " + plugin.Entrypoint + " " + plugin.Runtime.String())
	score := 0
	for _, word := range strings.Fields(intent) {
		word = strings.Trim(word, ".,:;()[]{}")
		if len(word) < 3 {
			continue
		}
		if strings.Contains(value, word) {
			score += 2
		}
	}
	if strings.Contains(intent, plugin.Runtime.String()) {
		score += 3
	}
	return score
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
