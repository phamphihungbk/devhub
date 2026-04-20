package projectrepo

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
)

func (r *projectRepositoryImpl) CreateOne(ctx context.Context, input *entity.Project) (project *entity.Project, err error) {
	const errLocation = "[repository project/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	projectsTable := table.Projects
	// SQL statement
	stmt := projectsTable.INSERT(
		projectsTable.AllColumns.Except(projectsTable.DefaultColumns), // Exclude columns with default values
	).MODEL(model.Projects{
		Name:        input.Name,
		Description: misc.ToPointer(input.Description),
		Status:      input.Status.String(),
		TeamID:      input.TeamID,
		ScmProvider: input.ScmProvider,
		CreatedBy:   input.CreatedBy,
		Environments: func() []string {
			envs := make([]string, 0, len(input.Environments))

			for _, env := range input.Environments {
				envs = append(envs, env.String())
			}

			return envs
		}(),
	}).RETURNING(projectsTable.AllColumns)
	query, args := stmt.Sql()

	var model Project
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating project", err.Error()))
	}

	project = model.ToEntity()
	if project == nil {
		return nil, errs.NewInternalServerError("failed to convert project model to entity", nil)
	}

	return project, nil
}
