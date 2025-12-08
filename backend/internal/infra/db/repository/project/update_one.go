package projectrepo

import (
	"context"
	"database/sql"
	"errors"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"

	postgres "github.com/go-jet/jet/v2/postgres"
)

func (r *projectRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateProjectInput) (project *entity.Project, err error) {
	const errLocation = "[repository project/update_one UpdateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	projectsTable := table.Projects
	var updateModel Project
	columns := make(postgres.ColumnList, 0)

	// build the update model
	if input.Name != nil {
		updateModel.Name = string(*input.Name)
		columns = append(columns, projectsTable.Name)
	}
	if input.Role != nil {
		updateModel.Role = string(*input.Role)
		columns = append(columns, projectsTable.Role)
	}

	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("no fields provided to update", nil)
	}

	// SQL statement
	stmt := projectsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(projectsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(projectsTable.AllColumns)
	query, args := stmt.Sql()

	var model Project
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("project not found", nil)
		}
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while updating user", err.Error()))
	}

	project = model.ToEntity()
	if project == nil {
		return nil, errs.NewInternalServerError("failed to convert project model to entity", nil)
	}

	return project, nil
}
