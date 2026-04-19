package releaserepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *releaseRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllReleasesFilter) (releases *entity.Releases, err error) {
	const errLocation = "[repository release/find_all FindAll] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	releasesTable := table.Releases
	stmt := postgres.SELECT(
		releasesTable.AllColumns,
	).
		FROM(releasesTable).
		WHERE(releasesTable.ServiceID.EQ(postgres.UUID(filter.ServiceID))).
		ORDER_BY(releasesTable.CreatedAt.DESC(), releasesTable.ID.DESC())

	query, args := stmt.Sql()

	var models Releases
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying releases", err.Error()))
	}

	releases = models.ToEntities()
	return releases, nil
}
