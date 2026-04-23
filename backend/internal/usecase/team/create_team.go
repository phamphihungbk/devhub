package usecase

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"
)

type CreateTeamInput struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	OwnerContact string `json:"owner_contact" validate:"required,min=3,max=255"`
}

func (u *teamUsecase) CreateTeam(ctx context.Context, input CreateTeamInput) (team *entity.Team, err error) {
	const errLocation = "[usecase team/create_team CreateTeam] "
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

	team = &entity.Team{
		Name:         input.Name,
		OwnerContact: input.OwnerContact,
	}

	created, err := u.teamRepository.CreateOne(ctx, team)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create team", nil))
	}

	return created, nil
}
