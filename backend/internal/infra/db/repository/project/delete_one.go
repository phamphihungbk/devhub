package projectrepo

import (
	"context"
	"database/sql"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
	"errors"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *projectRepositoryImpl) DeleteOne(ctx context.Context, id uuid.UUID) (project *entity.Project, err error) {
	const errLocation = "[repository project/delete Delete] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	projectsTable := table.Projects
	// SQL statement
	stmt := projectsTable.DELETE().
		WHERE(projectsTable.ID.EQ(postgres.UUID(id))).
		RETURNING(projectsTable.AllColumns)
	query, args := stmt.Sql()

	var model Project
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("project not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while deleting project", err.Error()))
	}

	project = model.ToEntity()
	if project == nil {
		return nil, errs.NewInternalServerError("failed to convert project model to entity", nil)
	}

	return project, nil
}
