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

type FindAllScaffoldRequestsInput struct {
	StartDate *time.Time        `json:"start_date" validate:"omitempty"`
	EndDate   *time.Time        `json:"end_date" validate:"omitempty,gtfield=StartDate"`
	Venue     *string           `json:"venue" validate:"omitempty,gt=0"`
	Limit     *int64            `json:"limit" validate:"required,gte=1,lte=100"`
	Offset    *int64            `json:"offset" validate:"required,gte=0"`
	SortBy    *string           `json:"sort_by" validate:"required_with=SortOrder,omitempty,oneof=date name venue"`
	SortOrder *entity.SortOrder `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	ProjectID string            `json:"project_id" validate:"required,uuid"`
}

func (u *scaffoldRequestUsecase) FindAllScaffoldRequests(ctx context.Context, input FindAllScaffoldRequestsInput) (scaffoldRequests entity.Page[entity.ScaffoldRequest], err error) {
	const errLocation = "[usecase scaffold_request/find_all_scaffold_request FindAllScaffoldRequests] "
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

	return entity.NewPage(u.findAllScaffoldRequests(ctx, input))
}

func (u *scaffoldRequestUsecase) findAllScaffoldRequests(ctx context.Context, input FindAllScaffoldRequestsInput) entity.PageProvider[entity.ScaffoldRequest] {
	return func() ([]entity.ScaffoldRequest, entity.PageProvider[entity.ScaffoldRequest], entity.Pagination, error) {
		// Fetch all scaffold requests with optional filters
		scaffoldRequests, count, err := u.scaffoldRequestRepository.FindAll(ctx, repository.FindAllScaffoldRequestsFilter{
			ProjectID: uuid.MustParse(input.ProjectID),
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    input.Offset,
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		})

		if err != nil {
			return entity.ScaffoldRequests{}, nil, entity.Pagination{}, errs.NewInternalServerError("failed to fetch scaffold requests", nil)
		}

		if scaffoldRequests == nil || len(misc.GetValue(scaffoldRequests)) == 0 {
			return entity.ScaffoldRequests{}, nil, entity.Pagination{}, nil
		}

		// Create pagination and next search criteria
		pagination := entity.NewPagination(count, misc.GetValue(input.Limit), misc.GetValue(input.Offset))
		nextSearchCriteria := FindAllScaffoldRequestsInput{
			ProjectID: input.ProjectID,
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    misc.ToPointer((*input.Limit) + (*input.Offset)),
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		}
		return misc.GetValue(scaffoldRequests), u.findAllScaffoldRequests(ctx, nextSearchCriteria), pagination, nil
	}
}
