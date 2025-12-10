package projectrepo

import (
	"context"
	"database/sql"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
	"errors"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *tokenRefreshRepositoryImpl) DeleteOne(ctx context.Context, id uuid.UUID) (token *entity.RefreshToken, err error) {
	const errLocation = "[repository token/delete_one DeleteOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	tokensTable := table.RefreshTokens
	// SQL statement
	stmt := tokensTable.DELETE().
		WHERE(tokensTable.ID.EQ(postgres.UUID(id))).
		RETURNING(tokensTable.AllColumns)
	query, args := stmt.Sql()

	var model RefreshToken
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("token not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while deleting token", err.Error()))
	}

	token = model.ToEntity()
	if token == nil {
		return nil, errs.NewInternalServerError("failed to convert token model to entity", nil)
	}

	return token, nil
}
