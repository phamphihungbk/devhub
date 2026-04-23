package teamrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	repository "devhub-backend/internal/domain/repository"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *teamRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateTeamInput) (_ *entity.Team, err error) {
	const errLocation = "[repository team/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	teamsTable := table.Teams
	var updateModel Team
	columns := make(postgres.ColumnList, 0)

	// build the update model
	if input.Name != nil {
		updateModel.Name = misc.GetValue(input.Name)
		columns = append(columns, teamsTable.Name)
	}

	if input.OwnerContact != nil {
		updateModel.OwnerContact = misc.GetValue(input.OwnerContact)
		columns = append(columns, teamsTable.OwnerContact)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	// SQL statement
	stmt := teamsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(teamsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(teamsTable.AllColumns)

	query, args := stmt.Sql()

	var model Team
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("team not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating team", err.Error()))
	}

	team := model.ToEntity()
	if team == nil {
		return nil, errs.NewInternalServerError("failed to convert team model to entity", nil)
	}

	return team, nil
}
