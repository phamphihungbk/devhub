package teamrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *teamRepositoryImpl) CreateOne(ctx context.Context, input *entity.Team) (team *entity.Team, err error) {
	const errLocation = "[repository team/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	teamsTable := table.Teams
	stmt := teamsTable.INSERT(
		teamsTable.AllColumns.Except(teamsTable.DefaultColumns),
	).MODEL(model.Teams{
		Name:         input.Name,
		OwnerContact: input.OwnerContact,
	}).RETURNING(teamsTable.AllColumns)

	query, args := stmt.Sql()

	var model Team
	if err := r.execer.GetContext(ctx, &model, query, args...); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating team", err.Error()))
	}

	team = model.ToEntity()
	if team == nil {
		return nil, errs.NewInternalServerError("failed to convert team model to entity", nil)
	}

	return team, nil
}
