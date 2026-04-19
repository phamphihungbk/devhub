package service

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	repository "devhub-backend/internal/domain/repository"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *serviceRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllServicesFilter) (services *entity.Services, total int64, err error) {
	const errLocation = "[repository service/find_all FindAll] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Build WHERE conditions for filtering
	whereClauses := []postgres.BoolExpression{}
	whereClauses = append(whereClauses, table.Services.ProjectID.EQ((postgres.UUID(filter.ProjectID))))

	if filter.StartDate != nil {
		whereClauses = append(whereClauses, table.Services.CreatedAt.GT_EQ(postgres.TimestampT(*filter.StartDate)))
	}
	if filter.EndDate != nil {
		whereClauses = append(whereClauses, table.Services.CreatedAt.LT_EQ(postgres.TimestampT(*filter.EndDate)))
	}

	// Get total count of services matching the filter
	countStmt := postgres.SELECT(
		postgres.COUNT(table.Services.ID).AS("total"),
	).FROM(table.Services)

	if len(whereClauses) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(whereClauses...))
	}

	countQuery, countArgs := countStmt.Sql()

	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while counting services", err.Error()))
	}

	// Get services with the same filter
	stmt := postgres.SELECT(
		table.Services.AllColumns,
	).FROM(table.Services)

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
				stmt = stmt.ORDER_BY(table.Services.Name.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Services.Name.ASC())
			}
		case "date":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Services.CreatedAt.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Services.CreatedAt.ASC())
			}
		}
	}

	query, args := stmt.Sql()

	var models Services
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while querying services", err.Error()))
	}

	services = models.ToEntities()
	return services, total, nil
}
