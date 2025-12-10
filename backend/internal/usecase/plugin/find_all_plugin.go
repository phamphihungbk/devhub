package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"time"

	"devhub-backend/internal/domain/errs"
	repository "devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"
)

type FindAllPluginsInput struct {
	StartDate *time.Time        `json:"start_date" validate:"omitempty"`
	EndDate   *time.Time        `json:"end_date" validate:"omitempty,gtfield=StartDate"`
	Venue     *string           `json:"venue" validate:"omitempty,gt=0"`
	Limit     *int64            `json:"limit" validate:"required,gte=1,lte=100"`
	Offset    *int64            `json:"offset" validate:"required,gte=0"`
	SortBy    *string           `json:"sort_by" validate:"required_with=SortOrder,omitempty,oneof=date name venue"`
	SortOrder *entity.SortOrder `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (u *pluginUsecase) FindAllPlugins(ctx context.Context, input FindAllPluginsInput) (plugins entity.Page[entity.Plugin], err error) {
	const errLocation = "[usecase plugin/find_all_plugin FindAllPlugins] "
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

	return entity.NewPage(u.findAllPlugins(ctx, input))
}

func (u *pluginUsecase) findAllPlugins(ctx context.Context, input FindAllPluginsInput) entity.PageProvider[entity.Plugin] {
	return func() ([]entity.Plugin, entity.PageProvider[entity.Plugin], entity.Pagination, error) {
		// Fetch all plugins with optional filters
		plugins, count, err := u.pluginRepository.FindAll(ctx, repository.FindAllPluginsFilter{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    input.Offset,
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		})

		if err != nil {
			return entity.Plugins{}, nil, entity.Pagination{}, errs.NewInternalServerError("failed to fetch plugins", nil)
		}

		if plugins == nil || len(misc.GetValue(plugins)) == 0 {
			return entity.Plugins{}, nil, entity.Pagination{}, nil
		}

		// Create pagination and next search criteria
		pagination := entity.NewPagination(count, misc.GetValue(input.Limit), misc.GetValue(input.Offset))
		nextSearchCriteria := FindAllPluginsInput{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    misc.ToPointer((*input.Limit) + (*input.Offset)),
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		}
		return misc.GetValue(plugins), u.findAllPlugins(ctx, nextSearchCriteria), pagination, nil
	}
}
