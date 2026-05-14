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
	"github.com/google/uuid"
)

func (r *jobRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (job *entity.Job, err error) {
	const errLocation = "[repository job/find_one FindOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	jobsTable := table.Jobs
	stmt := postgres.SELECT(
		jobsTable.AllColumns,
	).FROM(jobsTable).
		WHERE(jobsTable.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model Job
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("job not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying job by id", err.Error()))
	}

	job = model.ToEntity()
	if job == nil {
		return nil, errs.NewInternalServerError("failed to convert job model to entity", nil)
	}

	return job, nil
}
