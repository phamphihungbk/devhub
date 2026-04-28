package service

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *serviceRepositoryImpl) FindDependencies(ctx context.Context, serviceID uuid.UUID) (dependencies *entity.ServiceDependencies, err error) {
	const errLocation = "[repository service/find_dependencies FindDependencies] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	serviceDependenciesTable := table.ServiceDependencies
	dependsOnServicesTable := table.Services.AS("depends_on_services")
	stmt := postgres.SELECT(
		serviceDependenciesTable.AllColumns,
		dependsOnServicesTable.ProjectID.AS("depends_on_project_id"),
		dependsOnServicesTable.Name.AS("depends_on_name"),
		dependsOnServicesTable.RepoURL.AS("depends_on_repo_url"),
	).
		FROM(
			serviceDependenciesTable.INNER_JOIN(
				dependsOnServicesTable,
				serviceDependenciesTable.DependsOnServiceID.EQ(dependsOnServicesTable.ID),
			),
		).
		WHERE(serviceDependenciesTable.ServiceID.EQ(postgres.UUID(serviceID))).
		ORDER_BY(serviceDependenciesTable.CreatedAt.DESC())
	query, args := stmt.Sql()

	var models ServiceDependenciesWithService
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying service dependencies", err.Error()))
	}

	return models.ToEntities(), nil
}
