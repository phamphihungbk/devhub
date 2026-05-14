package jobrepo

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

func (r *jobRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateJobInput) (job *entity.Job, err error) {
	const errLocation = "[repository job/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	jobsTable := table.Jobs
	updateModel := model.Jobs{}
	columns := make(postgres.ColumnList, 0)

	if input.Status != nil {
		updateModel.Status = input.Status.String()
		columns = append(columns, jobsTable.Status)
	}
	if input.Result != nil {
		updateModel.Result = *input.Result
		columns = append(columns, jobsTable.Result)
	}
	if input.Error != nil {
		updateModel.Error = input.Error
		columns = append(columns, jobsTable.Error)
	}
	if input.Attempts != nil {
		updateModel.Attempts = int32(*input.Attempts)
		columns = append(columns, jobsTable.Attempts)
	}
	if input.StartedAt != nil {
		updateModel.StartedAt = input.StartedAt
		columns = append(columns, jobsTable.StartedAt)
	}
	if input.FinishedAt != nil {
		updateModel.FinishedAt = input.FinishedAt
		columns = append(columns, jobsTable.FinishedAt)
	}
	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}
	updateModel.UpdatedAt = time.Now()
	columns = append(columns, jobsTable.UpdatedAt)

	stmt := jobsTable.UPDATE(columns).
		MODEL(updateModel).
		WHERE(jobsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(jobsTable.AllColumns)

	query, args := stmt.Sql()

	var model Job
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("job not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating job", err.Error()))
	}

	job = model.ToEntity()
	if job == nil {
		return nil, errs.NewInternalServerError("failed to convert job model to entity", nil)
	}

	return job, nil
}
