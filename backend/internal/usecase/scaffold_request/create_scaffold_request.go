package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateScaffoldRequestInput struct {
	ProjectID   string            `json:"project_id" validate:"required,uuid"`
	Template    string            `json:"template" validate:"required,min=2,max=100"`
	Environment string            `json:"environment" validate:"required,oneof=dev staging prod"`
	Variables   ScaffoldVariables `json:"variables" validate:"required,dive"`
}

type ScaffoldVariables struct {
	ServiceName   string `json:"service_name" validate:"required"`
	Port          int    `json:"port" validate:"required"`
	Database      string `json:"database" validate:"required"`
	EnableLogging bool   `json:"enable_logging" validate:"required"`
}

func (u *scaffoldRequestUsecase) CreateScaffoldRequest(ctx context.Context, input CreateScaffoldRequestInput) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[usecase scaffold_request/create_scaffold_request CreateScaffoldRequest] "
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

	projectID, err := uuid.Parse(input.ProjectID)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid project ID", nil))
	}

	scaffoldRequest = &entity.ScaffoldRequest{
		ProjectID:   projectID,
		Template:    input.Template,
		Environment: new(entity.ProjectEnvironment).MustParse(input.Environment),
		Variables:   entity.ScaffoldRequestVariables(input.Variables),
	}
	created, err := u.scaffoldRequestRepository.CreateOne(ctx, scaffoldRequest)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create scaffold request", nil))
	}

	return created, nil
}
