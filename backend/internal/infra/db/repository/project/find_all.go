package projectrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	repository "devhub-backend/internal/domain/repository"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *projectRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllProjectsFilter) (projects *entity.Projects, total int64, err error) {
	const errLocation = "[repository project/find_all FindAll] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Build WHERE conditions for filtering
	whereClauses := []postgres.BoolExpression{}
	if filter.StartDate != nil {
		whereClauses = append(whereClauses, table.Projects.CreatedAt.GT_EQ(postgres.TimestampT(*filter.StartDate)))
	}
	if filter.EndDate != nil {
		whereClauses = append(whereClauses, table.Projects.CreatedAt.LT_EQ(postgres.TimestampT(*filter.EndDate)))
	}

	// Get total count of projects matching the filter
	countStmt := postgres.SELECT(
		postgres.COUNT(table.Projects.ID).AS("total"),
	).FROM(table.Projects)

	if len(whereClauses) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(whereClauses...))
	}

	countQuery, countArgs := countStmt.Sql()

	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while counting projects", err.Error()))
	}

	// Get projects with the same filter
	stmt := postgres.SELECT(
		table.Projects.AllColumns,
	).FROM(table.Projects)

	if len(whereClauses) > 0 {
		stmt = stmt.WHERE(postgres.AND(whereClauses...))
	}
	// Apply pagination
	if filter.Limit != nil {
		stmt = stmt.LIMIT(*filter.Limit)
	}
	if filter.Offset != nil {
		stmt = stmt.OFFSET(*filter.Offset)
	}
	// Apply sorting
	if filter.SortBy != nil {
		if filter.SortOrder == nil {
			filter.SortOrder = misc.ToPointer(entity.SortOrderAsc) // Default to ascending if not provided
		}
		switch *filter.SortBy {
		case "name":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Projects.Name.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Projects.Name.ASC())
			}
		case "date":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Projects.CreatedAt.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Projects.CreatedAt.ASC())
			}
		}
	}

	query, args := stmt.Sql()

	var models Projects
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while querying projects", err.Error()))
	}

	projects = models.ToEntities()
	return projects, total, nil
}
