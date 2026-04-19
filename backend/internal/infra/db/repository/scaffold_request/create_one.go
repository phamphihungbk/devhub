package scaffoldrequestrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *scaffoldRequestRepositoryImpl) CreateOne(ctx context.Context, input *entity.ScaffoldRequest) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[repository scaffold_request/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	scaffoldRequestsTable := table.ScaffoldRequests
	// SQL statement
	stmt := scaffoldRequestsTable.INSERT(
		scaffoldRequestsTable.AllColumns.Except(scaffoldRequestsTable.DefaultColumns),
	).MODEL(model.ScaffoldRequests{
		PluginID:    input.PluginID,
		RequestedBy: input.RequestedBy,
		Status:      input.Status.String(),
		ProjectID:   input.ProjectID,
		Environment: input.Environment.String(),
		Variables:   input.Variables.String(),
	}).RETURNING(scaffoldRequestsTable.AllColumns)
	query, args := stmt.Sql()

	var model ScaffoldRequest
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating scaffold request", err.Error()))
	}

	scaffoldRequest = model.ToEntity()
	if scaffoldRequest == nil {
		return nil, errs.NewInternalServerError("failed to convert scaffold request model to entity", nil)
	}

	return scaffoldRequest, nil
}
