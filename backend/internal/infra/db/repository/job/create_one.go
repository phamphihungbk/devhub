package jobrepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *jobRepositoryImpl) CreateOne(ctx context.Context, input *entity.Job) (job *entity.Job, err error) {
	const errLocation = "[repository job/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	jobsTable := table.Jobs
	stmt := jobsTable.INSERT(
		jobsTable.AllColumns.Except(jobsTable.DefaultColumns),
	).MODEL(model.Jobs{
		Type:         input.Type.String(),
		Status:       input.Status.String(),
		ResourceType: input.ResourceType.String(),
		ResourceID:   input.ResourceID,
		PluginID:     input.PluginID,
		Payload:      input.Payload.String(),
		CreatedBy:    input.CreatedBy,
	}).RETURNING(jobsTable.AllColumns)

	query, args := stmt.Sql()

	var model Job
	if err := r.execer.GetContext(ctx, &model, query, args...); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating job", err.Error()))
	}

	job = model.ToEntity()
	if job == nil {
		return nil, errs.NewInternalServerError("failed to convert job model to entity", nil)
	}

	return job, nil
}
