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

type FindAllServicesInput struct {
	ProjectID string `json:"project_id" validate:"required,uuid"`
}

func (u *serviceUsecase) FindAllServices(ctx context.Context, input FindAllServicesInput) (services entity.Services, err error) {
	const errLocation = "[usecase service/find_all_service FindAllServices] "
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

	found, _, err := u.serviceRepository.FindAll(ctx, repository.FindAllServicesFilter{
		ProjectID: uuid.MustParse(input.ProjectID),
	})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to fetch services", nil))
	}

	if found == nil {
		return entity.Services{}, nil
	}

	return misc.GetValue(found), nil
}
