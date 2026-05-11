package usecase

import (
	"context"
	"net/url"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"
)

type CreateReleaseFromCIInput struct {
	RepoURL     string `json:"repo_url" validate:"omitempty,url"`
	RepoOwner   string `json:"repo_owner" validate:"omitempty"`
	RepoName    string `json:"repo_name" validate:"omitempty"`
	Tag         string `json:"tag" validate:"required,git_revision,startswith=v"`
	Target      string `json:"target" validate:"omitempty,git_revision"`
	Name        string `json:"name" validate:"omitempty,max=255"`
	Notes       string `json:"notes" validate:"omitempty,max=5000"`
	HTMLURL     string `json:"html_url" validate:"omitempty,url"`
	Image       string `json:"image" validate:"omitempty,max=500"`
	ExternalRef string `json:"external_ref" validate:"omitempty,max=500"`
}

func (u *releaseUsecase) CreateReleaseFromCI(ctx context.Context, input CreateReleaseFromCIInput) (release *entity.Release, err error) {
	const errLocation = "[usecase release/create_release_from_ci CreateReleaseFromCI] "
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

	owner := strings.TrimSpace(input.RepoOwner)
	repoName := strings.TrimSpace(input.RepoName)
	if owner == "" || repoName == "" {
		owner, repoName, err = repoCoordinatesFromURL(input.RepoURL)
		if err != nil {
			return nil, misc.WrapError(err, errs.NewBadRequestError("repo_url or repo_owner/repo_name is required", nil))
		}
	}

	service, err := u.serviceRepository.FindOneByRepoCoordinates(ctx, owner, repoName)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("service not found for repository", map[string]string{"repo_owner": owner, "repo_name": repoName}))
	}

	externalRef := strings.TrimSpace(input.ExternalRef)
	if externalRef == "" {
		externalRef = strings.TrimSpace(input.Image)
	}

	release = &entity.Release{
		ServiceID:   service.ID,
		Tag:         input.Tag,
		Target:      input.Target,
		Name:        firstNonEmpty(input.Name, input.Tag),
		Notes:       input.Notes,
		HTMLURL:     input.HTMLURL,
		ExternalRef: externalRef,
		Status:      entity.ReleaseStatusCompleted,
		TriggeredBy: service.CreatedBy,
	}

	created, err := u.releaseRepository.CreateOne(ctx, release)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create release from ci webhook", nil))
	}

	return created, nil
}

func repoCoordinatesFromURL(repoURL string) (string, string, error) {
	parsed, err := url.Parse(strings.TrimSpace(repoURL))
	if err != nil {
		return "", "", err
	}

	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(segments) < 2 {
		return "", "", errs.NewBadRequestError("invalid repo_url", nil)
	}

	owner := strings.TrimSpace(segments[len(segments)-2])
	repoName := strings.TrimSuffix(strings.TrimSpace(segments[len(segments)-1]), ".git")
	if owner == "" || repoName == "" {
		return "", "", errs.NewBadRequestError("invalid repo_url", nil)
	}

	return owner, repoName, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
