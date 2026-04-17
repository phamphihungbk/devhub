package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateReleaseInput struct {
	ProjectID   string `json:"project_id" validate:"required,uuid"`
	PluginID    string `json:"plugin_id" validate:"required,uuid"`
	Tag         string `json:"tag" validate:"required,git_revision,startswith=v"`
	Target      string `json:"target" validate:"omitempty,git_revision"`
	Name        string `json:"name" validate:"omitempty,max=255"`
	Notes       string `json:"notes" validate:"omitempty,max=5000"`
	TriggeredBy string `json:"triggered_by" validate:"required,uuid"`
}

func (u *releaseUsecase) CreateRelease(ctx context.Context, input CreateReleaseInput) (release *entity.Release, err error) {
	const errLocation = "[usecase deployment/create_release CreateRelease] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
		validator.WithCustomValidator(validator.GitRevisionValidator{}),
	)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid environment", nil))
	}

	release = &entity.Release{
		ProjectID:   uuid.MustParse(input.ProjectID),
		PluginID:    uuid.MustParse(input.PluginID),
		Tag:         input.Tag,
		Target:      input.Target,
		Name:        input.Name,
		Notes:       input.Notes,
		Status:      entity.ReleaseStatusPending,
		TriggeredBy: uuid.MustParse(input.TriggeredBy),
	}

	created, err := u.releaseRepository.CreateOne(ctx, release)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create release", nil))
	}

	return created, nil
}
