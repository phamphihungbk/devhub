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

type UpdateTeamInput struct {
	ID           string  `json:"id" validate:"required,uuid"`
	Name         *string `json:"name" validate:"omitempty,min=2,max=100"`
	OwnerContact *string `json:"owner_contact" validate:"omitempty,min=3,max=255"`
}

func (u *teamUsecase) UpdateTeam(ctx context.Context, input UpdateTeamInput) (team *entity.Team, err error) {
	const errLocation = "[usecase team/update_team UpdateTeam] "
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

	updated, err := u.teamRepository.UpdateOne(ctx, repository.UpdateTeamInput{
		ID:           uuid.MustParse(input.ID),
		Name:         input.Name,
		OwnerContact: input.OwnerContact,
	})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update team", nil))
	}

	return updated, nil
}
