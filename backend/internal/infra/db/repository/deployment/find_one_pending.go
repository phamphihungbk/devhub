package deploymentrepo

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

func (r *deploymentRepositoryImpl) FindOnePending(ctx context.Context) (deployment *entity.Deployment, err error) {
	const errLocation = "[repository deployment/find_one_pending FindOnePending] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	deploymentsTable := table.Deployments
	stmt := postgres.SELECT(
		deploymentsTable.AllColumns,
	).
		FROM(deploymentsTable).
		WHERE(deploymentsTable.Status.EQ(postgres.String(entity.DeploymentStatusPending.String()))).
		ORDER_BY(deploymentsTable.CreatedAt.ASC(), deploymentsTable.ID.ASC()).
		LIMIT(1)

	query, args := stmt.Sql()

	var model Deployment
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying pending deployment", err.Error()))
	}

	deployment = model.ToEntity()
	if deployment == nil {
		return nil, errs.NewInternalServerError("failed to convert pending deployment model to entity", nil)
	}

	return deployment, nil
}
