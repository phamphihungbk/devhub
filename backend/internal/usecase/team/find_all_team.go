package usecase

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	repository "devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"
)

type FindAllTeamsInput struct {
	Limit     *int64            `json:"limit" validate:"required,gte=1,lte=100"`
	Offset    *int64            `json:"offset" validate:"required,gte=0"`
	SortBy    *string           `json:"sort_by" validate:"required_with=SortOrder,omitempty,oneof=date name"`
	SortOrder *entity.SortOrder `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (u *teamUsecase) FindAllTeams(ctx context.Context, input FindAllTeamsInput) (teams entity.Page[entity.Team], err error) {
	const errLocation = "[usecase team/find_all_team FindAllTeams] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err := vInstance.Struct(input); err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	return entity.NewPage(u.findAllTeams(ctx, input))
}

func (u *teamUsecase) findAllTeams(ctx context.Context, input FindAllTeamsInput) entity.PageProvider[entity.Team] {
	return func() ([]entity.Team, entity.PageProvider[entity.Team], entity.Pagination, error) {
		teams, count, err := u.teamRepository.FindAll(ctx, repository.FindAllTeamsFilter{
			Limit:     input.Limit,
			Offset:    input.Offset,
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		})
		if err != nil {
			return entity.Teams{}, nil, entity.Pagination{}, errs.NewInternalServerError("failed to fetch teams", nil)
		}

		if teams == nil || len(misc.GetValue(teams)) == 0 {
			return entity.Teams{}, nil, entity.Pagination{}, nil
		}

		pagination := entity.NewPagination(count, misc.GetValue(input.Limit), misc.GetValue(input.Offset))
		nextSearchCriteria := FindAllTeamsInput{
			Limit:     input.Limit,
			Offset:    misc.ToPointer((*input.Limit) + (*input.Offset)),
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		}

		return misc.GetValue(teams), u.findAllTeams(ctx, nextSearchCriteria), pagination, nil
	}
}
