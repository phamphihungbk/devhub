package scaffoldrequestrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *scaffoldRequestRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateScaffoldRequestInput) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[repository scaffold_request/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	scaffoldRequestsTable := table.ScaffoldRequests
	updateModel := model.ScaffoldRequests{}
	columns := make(postgres.ColumnList, 0)

	if input.PluginID != nil {
		updateModel.PluginID = *input.PluginID
		columns = append(columns, scaffoldRequestsTable.PluginID)
	}
	if input.ProjectID != nil {
		updateModel.ProjectID = *input.ProjectID
		columns = append(columns, scaffoldRequestsTable.ProjectID)
	}
	if input.RequestedBy != nil {
		updateModel.RequestedBy = *input.RequestedBy
		columns = append(columns, scaffoldRequestsTable.RequestedBy)
	}
	if input.Status != nil {
		updateModel.Status = input.Status.String()
		columns = append(columns, scaffoldRequestsTable.Status)
	}
	if input.Environment != nil {
		updateModel.Environment = input.Environment.String()
		columns = append(columns, scaffoldRequestsTable.Environment)
	}
	if input.Variables != nil {
		updateModel.Variables = input.Variables.String()
		columns = append(columns, scaffoldRequestsTable.Variables)
	}
	if input.ApprovedBy != nil {
		updateModel.ApprovedBy = input.ApprovedBy
		columns = append(columns, scaffoldRequestsTable.ApprovedBy)
	}
	if input.ResultRepoURL != nil {
		updateModel.ResultRepoURL = input.ResultRepoURL
		columns = append(columns, scaffoldRequestsTable.ResultRepoURL)
	}
	if input.ApprovedAt != nil {
		updateModel.ApprovedAt = input.ApprovedAt
		columns = append(columns, scaffoldRequestsTable.ApprovedAt)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	updateModel.UpdatedAt = time.Now()
	columns = append(columns, scaffoldRequestsTable.UpdatedAt)

	stmt := scaffoldRequestsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(scaffoldRequestsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(scaffoldRequestsTable.AllColumns)

	query, args := stmt.Sql()

	var model ScaffoldRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("scaffold request not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating scaffold request", err.Error()))
	}

	scaffoldRequest = model.ToEntity()
	if scaffoldRequest == nil {
		return nil, errs.NewInternalServerError("failed to convert scaffold request model to entity", nil)
	}

	return scaffoldRequest, nil
}
