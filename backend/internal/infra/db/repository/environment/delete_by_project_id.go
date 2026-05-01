package environment

import (
	"context"

	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *environmentRepositoryImpl) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) (err error) {
	const errLocation = "[repository environment/delete_by_project_id DeleteByProjectID] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	environmentsTable := table.Environments
	stmt := environmentsTable.DELETE().
		WHERE(environmentsTable.ProjectID.EQ(postgres.UUID(projectID)))

	query, args := stmt.Sql()
	if _, err := r.execer.ExecContext(ctx, query, args...); err != nil {
		return misc.WrapError(err, errs.NewDatabaseError("error while deleting environments by project id", err.Error()))
	}

	return nil
}
