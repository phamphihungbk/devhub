package environment

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *environmentRepositoryImpl) FindAllByProjectID(ctx context.Context, projectID uuid.UUID) (environments *entity.Environments, err error) {
	const errLocation = "[repository environment/find_all_by_project_id FindAllByProjectID] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	environmentsTable := table.Environments
	stmt := postgres.SELECT(
		environmentsTable.AllColumns,
	).
		FROM(environmentsTable).
		WHERE(environmentsTable.ProjectID.EQ(postgres.UUID(projectID))).
		ORDER_BY(environmentsTable.CreatedAt.ASC(), environmentsTable.ID.ASC())

	query, args := stmt.Sql()

	var models Environments
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying environments by project id", err.Error()))
	}

	return models.ToEntities(), nil
}
