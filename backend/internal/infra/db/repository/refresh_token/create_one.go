package refreshtokenrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
)

func (r *tokenRefreshRepositoryImpl) CreateOne(ctx context.Context, input repository.CreateRefreshTokenInput) (token *entity.RefreshToken, err error) {
	const errLocation = "[repository token/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	tokensTable := table.RefreshTokens
	// SQL statement
	stmt := tokensTable.INSERT(
		tokensTable.AllColumns.Except(tokensTable.DefaultColumns), // Exclude columns with default values
	).MODEL(model.RefreshTokens{
		UserID:    input.UserID,
		Token:     input.Token,
		ExpiresAt: input.ExpiresAt,
	}).RETURNING(tokensTable.AllColumns)
	query, args := stmt.Sql()

	var model RefreshToken
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating token", err.Error()))
	}

	token = model.ToEntity()
	if token == nil {
		return nil, errs.NewInternalServerError("failed to convert token model to entity", nil)
	}

	return token, nil
}
