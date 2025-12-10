package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"

	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type RevokeTokenInput struct {
	UserID string `json:"user_id" validate:"required"`
}

func (u *authUsecase) RevokeToken(ctx context.Context, input RevokeTokenInput) (user *entity.RefreshToken, err error) {
	const errLocation = "[usecase auth/revoke_token RevokeToken] "
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

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid user ID", nil))
	}

	deleted, err := u.refreshTokenRepository.DeleteOne(ctx, userID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to delete refresh token", nil))
	}

	return deleted, nil
}
