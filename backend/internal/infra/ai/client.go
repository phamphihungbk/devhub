package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOllamaClient(baseURL string, model string, timeout time.Duration) *OllamaClient {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &OllamaClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *OllamaClient) PlanScaffold(ctx context.Context, input ScaffoldPlanningInput) (*ScaffoldPlan, error) {
	if strings.TrimSpace(input.Prompt) == "" {
		return nil, fmt.Errorf("prompt is required")
	}

	if len(input.Plugins) == 0 {
		return nil, fmt.Errorf("at least one scaffold plugin is required")
	}

	prompt, err := buildOllamaScaffoldPrompt(input)
	if err != nil {
		return nil, err
	}

	reqBody := map[string]any{
		"model":  c.model,
		"prompt": prompt,
		"stream": false,
		"format": "json",
		"options": map[string]any{
			"temperature": 0.1,
		},
	}

	raw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/generate",
		bytes.NewReader(raw),
	)
	if err != nil {
		return nil, fmt.Errorf("build ollama request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute ollama request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read ollama response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama request failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var ollamaResp struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
	}

	content := strings.TrimSpace(ollamaResp.Response)
	if content == "" {
		return nil, fmt.Errorf("ollama returned empty response")
	}

	var plan ScaffoldPlan
	if err := json.Unmarshal([]byte(content), &plan); err != nil {
		return nil, fmt.Errorf("decode scaffold plan json: %w; content=%s", err, content)
	}

	if strings.TrimSpace(plan.PluginName) == "" {
		return nil, fmt.Errorf("ai scaffold plan missing plugin_name")
	}

	if plan.Variables == nil {
		plan.Variables = map[string]any{}
	}

	return &plan, nil
}

func buildOllamaScaffoldPrompt(input ScaffoldPlanningInput) (string, error) {
	pluginsJSON, err := json.MarshalIndent(input.Plugins, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal plugin candidates: %w", err)
	}

	return fmt.Sprintf(`
Choose exactly one scaffold plugin and produce valid variables for that plugin.

Rules:
- Return JSON only.
- Do not invent plugin names.
- Do not invent variables outside the selected plugin schema.
- Use schema defaults when the user does not provide a value.
- service_name must be lowercase kebab-case.
- If port is not mentioned, use the schema default.
- If database is not mentioned, use the schema default.

Return this exact JSON shape:
{
  "plugin_name": "string",
  "variables": {},
  "confidence": 0.0,
  "reason": "string"
}

User prompt:
%s

Available plugins:
%s
`, input.Prompt, string(pluginsJSON)), nil
}
