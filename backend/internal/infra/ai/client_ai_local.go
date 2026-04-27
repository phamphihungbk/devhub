package ai

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

type LocalClient struct{}

type localPluginScore struct {
	plugin  PluginCandidate
	score   float64
	matches []string
}

var localTokenPattern = regexp.MustCompile(`[a-z0-9]+`)

var localStopTokens = map[string]struct{}{
	"a": {}, "an": {}, "and": {}, "are": {}, "build": {}, "can": {}, "create": {}, "for": {}, "from": {}, "generate": {},
	"i": {}, "in": {}, "is": {}, "me": {}, "new": {}, "of": {}, "please": {}, "scaffold": {}, "service": {}, "the": {},
	"to": {}, "use": {}, "want": {}, "with": {},
}

var localTokenAliases = map[string][]string{
	"api":        {"api", "http", "rest", "backend", "server"},
	"backend":    {"backend", "api", "http", "server"},
	"db":         {"db", "database", "postgres", "postgresql", "mysql", "mariadb", "mongodb"},
	"database":   {"database", "db", "postgres", "postgresql", "mysql", "mariadb", "mongodb"},
	"frontend":   {"frontend", "ui", "web", "react", "vue", "nextjs", "vite"},
	"go":         {"go", "golang"},
	"golang":     {"go", "golang"},
	"grpc":       {"grpc", "rpc", "protobuf", "proto"},
	"http":       {"http", "api", "rest", "server"},
	"job":        {"job", "worker", "queue", "background"},
	"js":         {"js", "javascript", "node", "nodejs"},
	"next":       {"next", "nextjs", "react", "frontend"},
	"node":       {"node", "nodejs", "javascript", "js"},
	"nodejs":     {"node", "nodejs", "javascript", "js"},
	"postgres":   {"postgres", "postgresql", "database", "db"},
	"postgresql": {"postgres", "postgresql", "database", "db"},
	"python":     {"python", "py", "fastapi"},
	"react":      {"react", "frontend", "ui", "vite"},
	"rest":       {"rest", "api", "http"},
	"ui":         {"ui", "frontend", "web"},
	"vue":        {"vue", "frontend", "ui", "vite"},
	"web":        {"web", "frontend", "ui"},
	"worker":     {"worker", "job", "queue", "background"},
}

func NewLocalClient() *LocalClient {
	return &LocalClient{}
}

func (c *LocalClient) PlanScaffold(ctx context.Context, input ScaffoldPlanningInput) (*ScaffoldPlan, error) {
	if strings.TrimSpace(input.Prompt) == "" {
		return nil, fmt.Errorf("prompt is required")
	}
	if len(input.Plugins) == 0 {
		return nil, fmt.Errorf("at least one scaffold plugin is required")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	promptTokens := expandLocalTokens(tokenizeLocalText(input.Prompt))
	if len(promptTokens) == 0 {
		return nil, fmt.Errorf("prompt does not contain enough searchable terms")
	}

	scores := make([]localPluginScore, 0, len(input.Plugins))
	for _, plugin := range input.Plugins {
		pluginTokens := tokenizeLocalPlugin(plugin)
		matches := intersectLocalTokens(promptTokens, pluginTokens)
		score := float64(len(matches)) / float64(len(promptTokens))

		score += pluginNameMatchBonus(promptTokens, plugin)

		scores = append(scores, localPluginScore{
			plugin:  plugin,
			score:   math.Min(score, 1),
			matches: matches,
		})
	}

	sort.SliceStable(scores, func(i, j int) bool {
		if scores[i].score == scores[j].score {
			return scores[i].plugin.Name < scores[j].plugin.Name
		}
		return scores[i].score > scores[j].score
	})

	winner := scores[0]
	return &ScaffoldPlan{
		PluginName: winner.plugin.Name,
		Variables:  map[string]any{},
		Confidence: winner.score,
		Reason:     fmt.Sprintf("Matched %d of %d prompt tokens against plugin metadata.", len(winner.matches), len(promptTokens)),
		Matches:    winner.matches,
	}, nil
}

func tokenizeLocalPlugin(plugin PluginCandidate) map[string]struct{} {
	values := []string{
		plugin.Name,
		plugin.Type,
		plugin.Entrypoint,
		plugin.Description,
		strings.Join(plugin.Keywords, " "),
		strings.Join(plugin.Capabilities, " "),
	}

	return expandLocalTokens(tokenizeLocalText(strings.Join(values, " ")))
}

func tokenizeLocalText(value string) map[string]struct{} {
	tokens := map[string]struct{}{}
	for _, token := range localTokenPattern.FindAllString(strings.ToLower(value), -1) {
		if len(token) < 2 {
			continue
		}
		if _, ok := localStopTokens[token]; ok {
			continue
		}
		tokens[token] = struct{}{}
	}
	return tokens
}

func expandLocalTokens(tokens map[string]struct{}) map[string]struct{} {
	expanded := map[string]struct{}{}
	for token := range tokens {
		expanded[token] = struct{}{}
		for _, alias := range localTokenAliases[token] {
			expanded[alias] = struct{}{}
		}
	}
	return expanded
}

func intersectLocalTokens(left map[string]struct{}, right map[string]struct{}) []string {
	matches := make([]string, 0)
	for token := range left {
		if _, ok := right[token]; ok {
			matches = append(matches, token)
		}
	}
	sort.Strings(matches)
	return matches
}

func tokenSetContains(tokens map[string]struct{}, token string) bool {
	_, ok := tokens[strings.TrimSpace(strings.ToLower(token))]
	return ok
}

func pluginNameMatchBonus(promptTokens map[string]struct{}, plugin PluginCandidate) float64 {
	bonus := 0.0
	for token := range tokenizeLocalText(plugin.Name) {
		if tokenSetContains(promptTokens, token) {
			bonus += 0.05
		}
	}
	if bonus > 0.15 {
		return 0.15
	}
	return bonus
}
