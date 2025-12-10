package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateProjectInput struct {
	Name         string   `json:"name" validate:"required,min=2,max=100"`
	Description  *string  `json:"description" validate:"min=0,max=500"`
	Environments []string `json:"environments" validate:"required,dive,required,oneof=prod dev staging"`
	CreatedBy    string   `json:"created_by" validate:"required,uuid"`
}

func (u *projectUsecase) CreateProject(ctx context.Context, input CreateProjectInput) (project *entity.Project, err error) {
	const errLocation = "[usecase project/create_project CreateProject] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	envs, err := new(entity.ProjectEnvironment).ParseList(input.Environments)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid environments", map[string]string{"details": err.Error()}))
	}

	project = &entity.Project{
		Environments: envs,
		Name:         input.Name,
		Description:  misc.GetValue(input.Description),
		CreatedBy:    uuid.MustParse(input.CreatedBy),
	}

	created, err := u.projectRepository.CreateOne(ctx, project)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create project", nil))
	}

	return created, nil
}
