package teamrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	repository "devhub-backend/internal/domain/repository"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *teamRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllTeamsFilter) (teams *entity.Teams, total int64, err error) {
	const errLocation = "[repository team/find_all FindAll] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	countStmt := postgres.SELECT(
		postgres.COUNT(table.Teams.ID).AS("total"),
	).FROM(table.Teams)

	countQuery, countArgs := countStmt.Sql()
	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while counting teams", err.Error()))
	}

	stmt := postgres.SELECT(
		table.Teams.AllColumns,
	).FROM(table.Teams)

	if filter.Limit != nil {
		stmt = stmt.LIMIT(*filter.Limit)
	}
	if filter.Offset != nil {
		stmt = stmt.OFFSET(*filter.Offset)
	}
	if filter.SortBy != nil {
		if filter.SortOrder == nil {
			filter.SortOrder = misc.ToPointer(entity.SortOrderAsc)
		}
		switch *filter.SortBy {
		case "name":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Teams.Name.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Teams.Name.ASC())
			}
		default:
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Teams.CreatedAt.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Teams.CreatedAt.ASC())
			}
		}
	}

	query, args := stmt.Sql()

	var models Teams
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while querying teams", err.Error()))
	}

	teams = models.ToEntities()
	return teams, total, nil
}
