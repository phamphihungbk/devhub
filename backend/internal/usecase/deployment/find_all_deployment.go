package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"time"

	"devhub-backend/internal/domain/errs"
	repository "devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type FindAllDeploymentsInput struct {
	ProjectID string            `json:"project_id" validate:"required,uuid"`
	StartDate *time.Time        `json:"start_date" validate:"omitempty"`
	EndDate   *time.Time        `json:"end_date" validate:"omitempty,gtfield=StartDate"`
	Venue     *string           `json:"venue" validate:"omitempty,gt=0"`
	Limit     *int64            `json:"limit" validate:"required,gte=1,lte=100"`
	Offset    *int64            `json:"offset" validate:"required,gte=0"`
	SortBy    *string           `json:"sort_by" validate:"required_with=SortOrder,omitempty,oneof=date name venue"`
	SortOrder *entity.SortOrder `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (u *deploymentUsecase) FindAllDeployments(ctx context.Context, input FindAllDeploymentsInput) (deployments entity.Page[entity.Deployment], err error) {
	const errLocation = "[usecase deployment/find_all_deployment FindAllDeployments] "
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

	return entity.NewPage(u.findAllDeployments(ctx, input))
}

func (u *deploymentUsecase) findAllDeployments(ctx context.Context, input FindAllDeploymentsInput) entity.PageProvider[entity.Deployment] {
	return func() ([]entity.Deployment, entity.PageProvider[entity.Deployment], entity.Pagination, error) {
		// Fetch all deployments with optional filters
		deployments, count, err := u.deploymentRepository.FindAll(ctx, repository.FindAllDeploymentsFilter{
			ProjectID: uuid.MustParse(input.ProjectID),
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    input.Offset,
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		})

		if err != nil {
			return entity.Deployments{}, nil, entity.Pagination{}, errs.NewInternalServerError("failed to fetch deployments", nil)
		}

		if deployments == nil || len(misc.GetValue(deployments)) == 0 {
			return entity.Deployments{}, nil, entity.Pagination{}, nil
		}

		// Create pagination and next search criteria
		pagination := entity.NewPagination(count, misc.GetValue(input.Limit), misc.GetValue(input.Offset))
		nextSearchCriteria := FindAllDeploymentsInput{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    misc.ToPointer((*input.Limit) + (*input.Offset)),
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		}
		return misc.GetValue(deployments), u.findAllDeployments(ctx, nextSearchCriteria), pagination, nil
	}
}
