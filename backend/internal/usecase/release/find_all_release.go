package usecase

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type FindAllReleasesInput struct {
	ServiceID string `json:"service_id" validate:"required,uuid"`
}

func (u *releaseUsecase) FindAllReleases(ctx context.Context, input FindAllReleasesInput) (releases entity.Releases, err error) {
	const errLocation = "[usecase release/find_all_release FindAllReleases] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err = vInstance.Struct(input); err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	found, err := u.releaseRepository.FindAll(ctx, repository.FindAllReleasesFilter{
		ServiceID: uuid.MustParse(input.ServiceID),
	})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to fetch releases", nil))
	}

	if found == nil {
		return entity.Releases{}, nil
	}

	return misc.GetValue(found), nil
}
