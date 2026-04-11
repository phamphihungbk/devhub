package refreshtokenrepo

import (
	"context"
	"database/sql"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
	"errors"
	"time"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *tokenRefreshRepositoryImpl) FindOne(ctx context.Context, input repository.FindOneRefreshTokenInput) (token *entity.RefreshToken, err error) {
	const errLocation = "[repository token/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	tokensTable := table.RefreshTokens
	// SQL statement
	stmt := postgres.SELECT(
		tokensTable.AllColumns,
	).
		FROM(table.RefreshTokens).
		WHERE(tokensTable.UserID.EQ(postgres.UUID(input.UserID))).
		WHERE(tokensTable.ExpiresAt.LT(postgres.TimestampT(time.Now())))

	query, args := stmt.Sql()

	var model RefreshToken
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("scaffold request not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying scaffold request by id", err.Error()))
	}

	token = model.ToEntity()
	if token == nil {
		return nil, errs.NewInternalServerError("failed to convert token model to entity", nil)
	}

	return token, nil
}
