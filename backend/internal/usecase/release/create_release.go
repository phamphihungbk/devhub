package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	releaseprovider "devhub-backend/internal/infra/release_provider"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type CreateReleaseInput struct {
	ProjectID   string `json:"project_id" validate:"required,uuid"`
	Tag         string `json:"tag" validate:"required,git_revision,startswith=v"`
	Target      string `json:"target" validate:"omitempty,git_revision"`
	Name        string `json:"name" validate:"omitempty,max=255"`
	Notes       string `json:"notes" validate:"omitempty,max=5000"`
	TriggeredBy string `json:"triggered_by" validate:"required,uuid"`
}

func (u *releaseUsecase) CreateRelease(ctx context.Context, input CreateReleaseInput) (release *entity.Release, err error) {
	const errLocation = "[usecase release/create_release CreateRelease] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
		validator.WithCustomValidator(validator.GitRevisionValidator{}),
	)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err := vInstance.Struct(input); err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	if u.projectRepository == nil {
		return nil, errs.NewInternalServerError("project repository is required", nil)
	}
	if u.releaseRepository == nil {
		return nil, errs.NewInternalServerError("release repository is required", nil)
	}
	if len(u.releaseClients) == 0 {
		return nil, errs.NewInternalServerError("release clients are required", nil)
	}

	project, err := u.projectRepository.FindOne(ctx, uuid.MustParse(input.ProjectID))
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errs.NewNotFoundError("project not found", nil)
	}

	provider := strings.ToLower(strings.TrimSpace(project.RepoProvider))
	releaseClient, ok := u.releaseClients[provider]
	if !ok {
		return nil, errs.NewBadRequestError("unsupported project repo_provider for release creation", map[string]string{"details": project.RepoProvider})
	}

	owner, repo, err := parseRepoOwnerAndName(project.RepoURL)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid project repo_url", map[string]string{"details": err.Error()}))
	}

	target := strings.TrimSpace(input.Target)
	if target == "" {
		target = "main"
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		name = strings.TrimSpace(input.Tag)
	}

	created, statusCode, responseBody, err := releaseClient.CreateRelease(ctx, releaseprovider.CreateReleaseInput{
		Owner:           owner,
		Repo:            repo,
		TagName:         strings.TrimSpace(input.Tag),
		TargetCommitish: target,
		Name:            name,
		Body:            strings.TrimSpace(input.Notes),
	})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewThirdPartyError("failed to create repository release", nil))
	}

	if created == nil {
		switch statusCode {
		case 409:
			return nil, errs.NewConflictError("release tag already exists", map[string]string{"details": responseBody})
		case 404:
			return nil, errs.NewNotFoundError("repository not found in release provider", map[string]string{"details": responseBody})
		default:
			return nil, errs.NewThirdPartyError("release provider creation failed", map[string]string{"details": responseBody})
		}
	}

	release = &entity.Release{
		ProjectID:   uuid.MustParse(input.ProjectID),
		Tag:         created.TagName,
		Target:      target,
		Name:        created.Name,
		Notes:       created.Body,
		HTMLURL:     created.HTMLURL,
		ExternalRef: strings.TrimSpace(created.ID),
		TriggeredBy: uuid.MustParse(input.TriggeredBy),
	}

	persistedRelease, err := u.releaseRepository.CreateOne(ctx, release)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to persist release history", nil))
	}

	return persistedRelease, nil
}

func parseRepoOwnerAndName(repoURL string) (string, string, error) {
	raw := strings.TrimSpace(repoURL)
	if raw == "" {
		return "", "", fmt.Errorf("repo_url is required")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", "", fmt.Errorf("parse repo url: %w", err)
	}

	path := strings.Trim(parsed.Path, "/")
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("repo url path must include owner and repository name")
	}

	owner := strings.TrimSpace(parts[len(parts)-2])
	repo := strings.TrimSpace(parts[len(parts)-1])
	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("repo url path must include owner and repository name")
	}

	return owner, repo, nil
}
