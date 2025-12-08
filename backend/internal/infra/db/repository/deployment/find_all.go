package userrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	repository "devhub-backend/internal/domain/repository"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *deploymentRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllDeploymentsFilter) (deployments *entity.Deployments, total int64, err error) {
	const errLocation = "[repository deployment/find_all FindAll] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Build WHERE conditions for filtering
	whereClauses := []postgres.BoolExpression{}
	if filter.StartDate != nil {
		whereClauses = append(whereClauses, table.Users.CreatedAt.GT_EQ(postgres.TimestampT(*filter.StartDate)))
	}
	if filter.EndDate != nil {
		whereClauses = append(whereClauses, table.Users.CreatedAt.LT_EQ(postgres.TimestampT(*filter.EndDate)))
	}

	// Get total count of users matching the filter
	countStmt := postgres.SELECT(
		postgres.COUNT(table.Users.ID).AS("total"),
	).FROM(table.Users)

	if len(whereClauses) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(whereClauses...))
	}

	countQuery, countArgs := countStmt.Sql()

	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while counting deployments", err.Error()))
	}

	// Get deployments with the same filter
	stmt := postgres.SELECT(
		table.Deployments.AllColumns,
	).FROM(table.Deployments)

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
				stmt = stmt.ORDER_BY(table.Deployments.Name.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Deployments.Name.ASC())
			}
		case "date":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Deployments.CreatedAt.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Deployments.CreatedAt.ASC())
			}
		}
	}

	query, args := stmt.Sql()

	var models Deployments
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while querying deployments", err.Error()))
	}

	deployments = models.ToEntities()
	return deployments, total, nil
}
