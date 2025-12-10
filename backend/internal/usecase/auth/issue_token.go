package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	"time"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	jwt "github.com/golang-jwt/jwt/v5"
)

type IssueTokenInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *authUsecase) IssueToken(ctx context.Context, input IssueTokenInput) (token *entity.Token, err error) {
	const errLocation = "[usecase auth/login_user IssueToken] "
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

	user, err := u.userRepository.FindOneByEmail(ctx, input.Email)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find user by email", nil))
	}

	userPasswordValid := misc.CheckPassword(user.PasswordHash, input.Password)

	if !userPasswordValid {
		return nil, misc.WrapError(err, errs.NewBadRequestError("password is incorrect", nil))
	}

	refreshToken, err := u.IssueRefreshToken(ctx, user)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to issue refresh token", nil))
	}

	accessToken, err := u.IssueAccessToken(ctx, user)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to issue access token", nil))
	}

	return &entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) IssueRefreshToken(ctx context.Context, user *entity.User) (refreshTokenString string, err error) {
	const errLocation = "[usecase auth/issue_token IssueRefreshToken] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	refreshTokenString, err = misc.GenerateSecureToken(64)

	if err != nil {
		return "", misc.WrapError(err, errs.NewInternalServerError("failed to generate refresh token", nil))
	}

	refreshToken, _ := u.refreshTokenRepository.FindOne(ctx, repository.FindOneRefreshTokenInput{
		UserID: user.ID,
	})

	if refreshToken != nil {
		return "", misc.WrapError(err, errs.NewInternalServerError("failed to generate new refresh token", nil))
	}

	created, err := u.refreshTokenRepository.CreateOne(ctx, repository.CreateRefreshTokenInput{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	})

	if err != nil {
		return "", misc.WrapError(err, errs.NewInternalServerError("failed to create refresh token", nil))
	}

	return created.Token, nil
}

func (u *authUsecase) IssueAccessToken(ctx context.Context, user *entity.User) (accessToken string, err error) {
	const errLocation = "[usecase auth/issue_token IssueAccessToken] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	now := time.Now()
	claims := (entity.AccessToken{
		UserID:    user.ID,
		Role:      user.Role,
		IssuedAt:  now,
		Issuer:    u.tokenConfig.Issuer,
		ExpiresAt: now.Add(time.Duration(u.tokenConfig.Duration) * time.Second),
	}).ToJWTClaims()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(u.tokenConfig.Secret))

	if err != nil {
		return "", misc.WrapError(err, errs.NewInternalServerError("failed to sign access token", nil))
	}

	return accessToken, nil
}
