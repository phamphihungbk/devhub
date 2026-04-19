package scaffoldrequestrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	repository "devhub-backend/internal/domain/repository"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *scaffoldRequestRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllScaffoldRequestsFilter) (scaffoldRequests *entity.ScaffoldRequests, total int64, err error) {
	const errLocation = "[repository scaffold_request/find_all FindAll] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Build WHERE conditions for filtering
	whereClauses := []postgres.BoolExpression{}

	whereClauses = append(whereClauses, table.ScaffoldRequests.ProjectID.EQ(postgres.UUID(filter.ProjectID)))

	// Get total count of users matching the filter
	countStmt := postgres.SELECT(
		postgres.COUNT(table.ScaffoldRequests.ID).AS("total"),
	).FROM(table.ScaffoldRequests)

	if len(whereClauses) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(whereClauses...))
	}

	countQuery, countArgs := countStmt.Sql()

	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while counting scaffold requests", err.Error()))
	}

	// Get scaffold requests with the same filter
	stmt := postgres.SELECT(
		table.ScaffoldRequests.AllColumns,
	).FROM(table.ScaffoldRequests)

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
				stmt = stmt.ORDER_BY(table.ScaffoldRequests.CreatedAt.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.ScaffoldRequests.CreatedAt.ASC())
			}
		case "date":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.ScaffoldRequests.CreatedAt.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.ScaffoldRequests.CreatedAt.ASC())
			}
		}
	}

	query, args := stmt.Sql()

	var models ScaffoldRequests
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while querying scaffold requests", err.Error()))
	}

	scaffoldRequests = models.ToEntities()
	return scaffoldRequests, total, nil
}
