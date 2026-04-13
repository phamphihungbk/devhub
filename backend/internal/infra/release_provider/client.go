package releaseprovider

import (
	"bytes"
	"context"
	"devhub-backend/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	Provider() string
	CreateRelease(ctx context.Context, input CreateReleaseInput) (*Release, int, string, error)
}

type GiteaClient struct {
	baseURL     string
	externalURL string
	token       string
	httpClient  *http.Client
}

type CreateReleaseInput struct {
	Owner           string
	Repo            string
	TagName         string
	TargetCommitish string
	Name            string
	Body            string
	Draft           bool
	Prerelease      bool
}

type Release struct {
	ID              string `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	HTMLURL         string `json:"html_url"`
	TarballURL      string `json:"tarball_url"`
	ZipballURL      string `json:"zipball_url"`
	Draft           bool   `json:"draft"`
	Prerelease      bool   `json:"prerelease"`
}

func NewClient(cfg config.GiteaConfig) Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &GiteaClient{
		baseURL:     strings.TrimRight(strings.TrimSpace(cfg.URL), "/"),
		externalURL: strings.TrimRight(strings.TrimSpace(cfg.ExternalURL), "/"),
		token:       strings.TrimSpace(cfg.Token),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *GiteaClient) Provider() string {
	return "gitea"
}

func (c *GiteaClient) CreateRelease(ctx context.Context, input CreateReleaseInput) (*Release, int, string, error) {
	payload := map[string]any{
		"tag_name":         input.TagName,
		"target_commitish": input.TargetCommitish,
		"name":             input.Name,
		"body":             input.Body,
		"draft":            input.Draft,
		"prerelease":       input.Prerelease,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, "", fmt.Errorf("marshal gitea release payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/repos/%s/%s/releases", c.baseURL, input.Owner, input.Repo),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, 0, "", fmt.Errorf("build gitea release request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, "", fmt.Errorf("execute gitea release request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, "", fmt.Errorf("read gitea release response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, strings.TrimSpace(string(respBody)), nil
	}

	var apiRelease struct {
		ID              int64  `json:"id"`
		TagName         string `json:"tag_name"`
		TargetCommitish string `json:"target_commitish"`
		Name            string `json:"name"`
		Body            string `json:"body"`
		HTMLURL         string `json:"html_url"`
		TarballURL      string `json:"tarball_url"`
		ZipballURL      string `json:"zipball_url"`
		Draft           bool   `json:"draft"`
		Prerelease      bool   `json:"prerelease"`
	}
	if err := json.Unmarshal(respBody, &apiRelease); err != nil {
		return nil, resp.StatusCode, strings.TrimSpace(string(respBody)), fmt.Errorf("decode gitea release response: %w", err)
	}

	release := &Release{
		ID:              fmt.Sprintf("%d", apiRelease.ID),
		TagName:         apiRelease.TagName,
		TargetCommitish: apiRelease.TargetCommitish,
		Name:            apiRelease.Name,
		Body:            apiRelease.Body,
		HTMLURL:         apiRelease.HTMLURL,
		TarballURL:      apiRelease.TarballURL,
		ZipballURL:      apiRelease.ZipballURL,
		Draft:           apiRelease.Draft,
		Prerelease:      apiRelease.Prerelease,
	}

	if strings.TrimSpace(release.HTMLURL) == "" && c.externalURL != "" {
		release.HTMLURL = fmt.Sprintf("%s/%s/%s/releases/tag/%s", c.externalURL, input.Owner, input.Repo, input.TagName)
	}

	return release, resp.StatusCode, strings.TrimSpace(string(respBody)), nil
}
