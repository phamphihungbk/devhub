package usecase

import (
	"context"
	"regexp"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

type SuggestScaffoldInput struct {
	ServiceID          string   `json:"service_id"`
	ServiceName        string   `json:"service_name"`
	ProjectName        string   `json:"project_name"`
	ProjectDescription string   `json:"project_description"`
	RepoURL            string   `json:"repo_url"`
	Environment        string   `json:"environment"`
	Environments       []string `json:"environments"`
}

type ScaffoldSuggestion struct {
	Source      string                          `json:"source"`
	Environment string                          `json:"environment"`
	Variables   entity.ScaffoldRequestVariables `json:"variables"`
	Rationale   []string                        `json:"rationale"`
}

var nonServiceNameCharacterPattern = regexp.MustCompile(`[^a-z0-9]+`)

func (u *serviceUsecase) SuggestScaffold(ctx context.Context, input SuggestScaffoldInput) (suggestion ScaffoldSuggestion, err error) {
	const errLocation = "[usecase service/suggest_scaffold SuggestScaffold] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	serviceName := normalizeScaffoldServiceName(firstNonEmpty(input.ServiceName, input.ProjectName, inferNameFromDescription(input.ProjectDescription), "new-service"))
	modulePath := inferScaffoldModulePath(input.RepoURL)
	environment := pickScaffoldEnvironment(input.Environment, input.Environments)
	database := suggestScaffoldDatabase(serviceName, modulePath, input.ProjectDescription)
	enableLogging := suggestScaffoldLogging(input.ProjectDescription)

	if environment == "" {
		return ScaffoldSuggestion{}, errs.NewBadRequestError("unable to infer scaffold environment", nil)
	}

	rationale := []string{
		"Service name was normalized from service, project, or description context.",
		"Port was selected from a stable service-name hash to avoid random churn.",
	}

	if strings.TrimSpace(input.ProjectDescription) != "" {
		rationale = append(rationale, "Project description was used to infer scaffold defaults.")
	}

	if modulePath != "" {
		rationale = append(rationale, "Module path was inferred from the repository URL.")
	} else {
		rationale = append(rationale, "Module path needs manual input because the repository URL is empty.")
	}

	if input.Environment != "" {
		rationale = append(rationale, "Environment follows the current scaffold form selection.")
	} else {
		rationale = append(rationale, "Environment uses the first project environment, falling back to dev.")
	}

	return ScaffoldSuggestion{
		Source:      "local-heuristic-v1",
		Environment: environment,
		Variables: entity.ScaffoldRequestVariables{
			ServiceName:   serviceName,
			ModulePath:    modulePath,
			Port:          suggestScaffoldPort(serviceName),
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

func normalizeScaffoldServiceName(value string) string {
	normalized := strings.Trim(nonServiceNameCharacterPattern.ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "-"), "-")
	if normalized == "" {
		return "new-service"
	}
	return normalized
}

func inferScaffoldModulePath(repoURL string) string {
	trimmed := strings.TrimSpace(strings.TrimSuffix(repoURL, ".git"))
	if trimmed == "" {
		return ""
	}

	trimmed = strings.TrimPrefix(trimmed, "https://")
	trimmed = strings.TrimPrefix(trimmed, "http://")
	trimmed = strings.TrimPrefix(trimmed, "git@")

	return strings.Replace(trimmed, ":", "/", 1)
}

func inferNameFromDescription(description string) string {
	words := strings.Fields(strings.ToLower(description))
	if len(words) == 0 {
		return ""
	}

	stopWords := map[string]struct{}{
		"a": {}, "an": {}, "and": {}, "for": {}, "of": {}, "the": {}, "to": {}, "with": {},
		"build": {}, "builds": {}, "service": {}, "services": {}, "system": {}, "platform": {},
	}

	selected := make([]string, 0, 3)
	for _, word := range words {
		word = strings.Trim(word, ".,:;()[]{}")
		if _, ok := stopWords[word]; ok || word == "" {
			continue
		}
		selected = append(selected, word)
		if len(selected) == 3 {
			break
		}
	}

	return strings.Join(selected, "-")
}

func pickScaffoldEnvironment(selected string, environments []string) string {
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

func suggestScaffoldPort(serviceName string) int {
	hash := 0
	for _, char := range serviceName {
		hash += int(char)
	}
	return 8000 + hash%1000
}

func suggestScaffoldDatabase(serviceName string, modulePath string, projectDescription string) string {
	value := strings.ToLower(serviceName + " " + modulePath + " " + projectDescription)
	switch {
	case strings.Contains(value, "cache"), strings.Contains(value, "redis"):
		return "redis"
	case strings.Contains(value, "mongo"), strings.Contains(value, "document"):
		return "mongodb"
	case strings.Contains(value, "analytics"), strings.Contains(value, "warehouse"), strings.Contains(value, "reporting"):
		return "postgres"
	default:
		return "postgres"
	}
}

func suggestScaffoldLogging(projectDescription string) bool {
	value := strings.ToLower(projectDescription)
	if strings.Contains(value, "disable logging") || strings.Contains(value, "no logging") {
		return false
	}
	return true
}
