package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type UpdateUserInput struct {
	ID     string  `json:"id" validate:"required,uuid"`
	Name   *string `json:"name" validate:"min=2,max=100"`
	Role   *string `json:"role" validate:"omitempty,oneof=platform_admin org_admin team_lead developer viewer"`
	TeamID *string `json:"team_id" validate:"omitempty,uuid"`
}

func (u *userUsecase) UpdateUser(ctx context.Context, input UpdateUserInput) (user *entity.User, err error) {
	const errLocation = "[usecase user/update_user UpdateUser] "
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

	updated, err := u.userRepository.UpdateOne(ctx, repository.UpdateUserInput{
		ID:   uuid.MustParse(input.ID),
		Name: input.Name,
		Role: func() *entity.UserRole {
			if input.Role == nil {
				return nil
			}
			role, err := new(entity.UserRole).Parse(*input.Role)
			if err != nil {
				return nil
			}
			return &role
		}(),
		TeamID: func() *uuid.UUID {
			if input.TeamID == nil {
				return nil
			}
			parsed := uuid.MustParse(*input.TeamID)
			return &parsed
		}(),
	})

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update user", nil))
	}

	return updated, nil
}
