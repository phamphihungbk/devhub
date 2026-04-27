package ai

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"devhub-backend/internal/domain/entity"
)

type LocalScaffoldSuggestionGenerator struct{}

var (
	scaffoldSuggestionServiceNameCharacterPattern = regexp.MustCompile(`[^a-z0-9]+`)
	scaffoldSuggestionPortPattern                 = regexp.MustCompile(`\bport\s*[:=]?\s*(\d{1,5})\b`)
)

var stopWords = map[string]struct{}{
	"a": {}, "an": {}, "and": {}, "api": {}, "app": {}, "build": {}, "create": {}, "for": {}, "generate": {},
	"database": {}, "db": {}, "dev": {}, "environment": {}, "go": {}, "http": {}, "logging": {}, "mariadb": {}, "mongodb": {},
	"mysql": {}, "node": {}, "of": {}, "please": {}, "port": {}, "postgres": {}, "prod": {}, "python": {},
	"redis": {}, "scaffold": {}, "service": {}, "staging": {}, "structured": {}, "the": {}, "to": {}, "with": {},
}

func NewLocalScaffoldSuggestionGenerator() *LocalScaffoldSuggestionGenerator {
	return &LocalScaffoldSuggestionGenerator{}
}

func (g *LocalScaffoldSuggestionGenerator) GenerateScaffoldSuggestion(ctx context.Context, input ScaffoldSuggestionInput) (*ScaffoldSuggestion, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	serviceName := normalizeScaffoldSuggestionServiceName(firstNonEmpty(inferNameFromPrompt(input.Prompt), input.Project.Name, inferNameFromPrompt(input.Project.Description), "new-service"))
	suggestedEnvironments := inferScaffoldSuggestionEnvironments(input.Prompt, input.ProjectEnvironments)
	environment := pickScaffoldSuggestionEnvironment("", suggestedEnvironments)
	modulePath := suggestScaffoldSuggestionModulePath(input.Plugin, input.Project.Name, serviceName)
	database := suggestScaffoldSuggestionDatabase(input.Prompt, input.Project.Description, input.Plugin)
	port := suggestScaffoldSuggestionPort(input.Prompt, serviceName)
	enableLogging := suggestScaffoldSuggestionLogging(input.Prompt, input.Project.Description)

	rationale := []string{
		"Prompt was analyzed by the local token ranking client.",
		fmt.Sprintf("Winning plugin score: %.0f%%.", input.Plan.Confidence*100),
		"Scaffolder plugin was selected from enabled scaffold_request plugins only.",
		"Service name was inferred from the user prompt and project context.",
		"Selected plugin: " + input.Plugin.Name + ".",
	}
	if len(input.Plan.Matches) > 0 {
		rationale = append(rationale, "Matched tokens: "+strings.Join(input.Plan.Matches, ", ")+".")
	}
	if hasExplicitScaffoldSuggestionPort(input.Prompt) {
		rationale = append(rationale, "Port was taken from the prompt.")
	} else {
		rationale = append(rationale, "Port was selected from a stable service-name hash.")
	}
	if len(suggestedEnvironments) > 0 {
		rationale = append(rationale, "Environments were inferred from the prompt or project settings.")
	}
	if strings.TrimSpace(input.Project.Description) != "" {
		rationale = append(rationale, "Project description was used as additional intent context.")
	}

	return &ScaffoldSuggestion{
		Source:      "local-token-ranker-v1",
		Environment: environment,
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
	normalized := strings.Trim(scaffoldSuggestionServiceNameCharacterPattern.ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "-"), "-")
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

func suggestScaffoldSuggestionModulePath(plugin entity.Plugin, projectName string, serviceName string) string {
	host := "gitea.devhub.local"
	if strings.Contains(plugin.Entrypoint, "github") {
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
	matches := scaffoldSuggestionPortPattern.FindStringSubmatch(strings.ToLower(prompt))
	if len(matches) != 2 {
		return 0, false
	}

	port, err := strconv.Atoi(matches[1])
	if err != nil || port < 1 || port > 65535 {
		return 0, false
	}

	return port, true
}

func hasExplicitScaffoldSuggestionPort(prompt string) bool {
	_, ok := inferScaffoldSuggestionPort(prompt)
	return ok
}

func suggestScaffoldSuggestionDatabase(prompt string, projectDescription string, plugin entity.Plugin) string {
	value := strings.ToLower(prompt + " " + projectDescription + " " + plugin.Description)
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
