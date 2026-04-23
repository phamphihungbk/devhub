package approvalrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	repository "devhub-backend/internal/domain/repository"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *approvalRepositoryImpl) FindAllApprovalRequests(ctx context.Context, filter repository.FindAllApprovalRequestsFilter) (requests *entity.ApprovalRequests, total int64, err error) {
	const errLocation = "[repository approval/find_all_requests FindAllApprovalRequests] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	approvalRequestsTable := table.ApprovalRequests
	whereClauses := []postgres.BoolExpression{}
	if filter.Status != nil {
		whereClauses = append(whereClauses, approvalRequestsTable.Status.EQ(postgres.String(filter.Status.String())))
	}

	countStmt := postgres.SELECT(
		postgres.COUNT(approvalRequestsTable.ID).AS("total"),
	).FROM(approvalRequestsTable)
	if len(whereClauses) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(whereClauses...))
	}

	countQuery, countArgs := countStmt.Sql()
	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while counting approval requests", err.Error()))
	}

	stmt := postgres.SELECT(
		approvalRequestsTable.AllColumns,
	).FROM(approvalRequestsTable)
	if len(whereClauses) > 0 {
		stmt = stmt.WHERE(postgres.AND(whereClauses...))
	}
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
		case "status":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(approvalRequestsTable.Status.DESC())
			} else {
				stmt = stmt.ORDER_BY(approvalRequestsTable.Status.ASC())
			}
		default:
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(approvalRequestsTable.CreatedAt.DESC())
			} else {
				stmt = stmt.ORDER_BY(approvalRequestsTable.CreatedAt.ASC())
			}
		}
	}

	query, args := stmt.Sql()

	var models ApprovalRequests
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, misc.WrapError(err, errs.NewDatabaseError("error while querying approval requests", err.Error()))
	}

	requests = models.ToEntities()
	return requests, total, nil
}
