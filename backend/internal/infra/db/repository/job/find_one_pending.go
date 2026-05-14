package jobrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *jobRepositoryImpl) FindOnePending(ctx context.Context, jobType entity.JobType) (job *entity.Job, err error) {
	const errLocation = "[repository job/find_one_pending FindOnePending] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	jobsTable := table.Jobs
	stmt := postgres.SELECT(
		jobsTable.AllColumns,
	).FROM(jobsTable).
		WHERE(
			jobsTable.Type.EQ(postgres.String(jobType.String())).
				AND(jobsTable.Status.EQ(postgres.String(entity.JobStatusQueued.String()))),
		).
		ORDER_BY(jobsTable.CreatedAt.ASC(), jobsTable.ID.ASC()).
		LIMIT(1)

	query, args := stmt.Sql()

	var model Job
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying pending job", err.Error()))
	}

	job = model.ToEntity()
	if job == nil {
		return nil, errs.NewInternalServerError("failed to convert job model to entity", nil)
	}

	return job, nil
}
