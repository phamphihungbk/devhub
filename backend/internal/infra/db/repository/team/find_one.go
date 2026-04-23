package teamrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *teamRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (_ *entity.Team, err error) {
	const errLocation = "[repository team/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	teamsTable := table.Teams
	stmt := postgres.SELECT(
		teamsTable.AllColumns,
	).
		FROM(table.Teams).
		WHERE(table.Teams.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model Team
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("team not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying team by id", err.Error()))
	}

	team := model.ToEntity()
	if team == nil {
		return nil, errs.NewInternalServerError("failed to convert team model to entity", nil)
	}

	return team, nil
}
