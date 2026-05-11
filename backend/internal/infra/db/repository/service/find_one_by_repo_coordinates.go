package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *serviceRepositoryImpl) FindOneByRepoCoordinates(ctx context.Context, owner string, repoName string) (service *entity.Service, err error) {
	const errLocation = "[repository service/find_one_by_repo_coordinates FindOneByRepoCoordinates] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	owner = strings.Trim(strings.TrimSpace(owner), "/")
	repoName = strings.Trim(strings.TrimSpace(repoName), "/")
	repoName = strings.TrimSuffix(repoName, ".git")
	if owner == "" || repoName == "" {
		return nil, errs.NewNotFoundError("service not found", nil)
	}

	repoPath := fmt.Sprintf("%%/%s/%s", owner, repoName)
	repoGitPath := repoPath + ".git"

	servicesTable := table.Services
	stmt := postgres.SELECT(
		servicesTable.AllColumns,
	).
		FROM(servicesTable).
		WHERE(
			servicesTable.RepoURL.LIKE(postgres.String(repoPath)).
				OR(servicesTable.RepoURL.LIKE(postgres.String(repoGitPath))),
		).
		LIMIT(1)

	query, args := stmt.Sql()

	var model Service
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("service not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while querying service by repo coordinates", err.Error()))
	}

	service = model.ToEntity()
	if service == nil {
		return nil, errs.NewInternalServerError("failed to convert service model to entity", nil)
	}

	return service, nil
}
