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

type FindAllUsersInput struct {
	StartDate *time.Time        `json:"start_date" validate:"omitempty"`
	EndDate   *time.Time        `json:"end_date" validate:"omitempty,gtfield=StartDate"`
	Limit     *int64            `json:"limit" validate:"required,gte=1,lte=100"`
	Offset    *int64            `json:"offset" validate:"required,gte=0"`
	SortBy    *string           `json:"sort_by" validate:"required_with=SortOrder,omitempty,oneof=date name email"`
	SortOrder *entity.SortOrder `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (u *userUsecase) FindAllUsers(ctx context.Context, input FindAllUsersInput) (users entity.Page[entity.User], err error) {
	const errLocation = "[usecase user/find_all_user FindAllUsers] "
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

	return entity.NewPage(u.findAllUsers(ctx, input))
}

func (u *userUsecase) findAllUsers(ctx context.Context, input FindAllUsersInput) entity.PageProvider[entity.User] {
	return func() ([]entity.User, entity.PageProvider[entity.User], entity.Pagination, error) {
		// Fetch all users with optional filters
		users, count, err := u.userRepository.FindAll(ctx, repository.FindAllUsersFilter{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    input.Offset,
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		})

		if err != nil {
			return entity.Users{}, nil, entity.Pagination{}, errs.NewInternalServerError("failed to fetch users", nil)
		}

		if users == nil || len(misc.GetValue(users)) == 0 {
			return entity.Users{}, nil, entity.Pagination{}, nil
		}

		// Create pagination and next search criteria
		pagination := entity.NewPagination(count, misc.GetValue(input.Limit), misc.GetValue(input.Offset))
		nextSearchCriteria := FindAllUsersInput{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Limit:     input.Limit,
			Offset:    misc.ToPointer((*input.Limit) + (*input.Offset)),
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		}
		return misc.GetValue(users), u.findAllUsers(ctx, nextSearchCriteria), pagination, nil
	}
}
