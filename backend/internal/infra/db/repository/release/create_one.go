package releaserepo

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	table "devhub-backend/internal/infra/db/model_gen/devhub/public/table"
	"devhub-backend/internal/util/misc"
)

func (r *releaseRepositoryImpl) CreateOne(ctx context.Context, input *entity.Release) (release *entity.Release, err error) {
	const errLocation = "[repository release/create_one CreateOne] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	releasesTable := table.Releases
	stmt := releasesTable.INSERT(
		releasesTable.AllColumns.Except(releasesTable.DefaultColumns),
	).MODEL(model.Releases{
		ServiceID:   input.ServiceID,
		PluginID:    input.PluginID,
		Tag:         input.Tag,
		Target:      input.Target,
		Name:        input.Name,
		Status:      input.Status.String(),
		Notes:       input.Notes,
		HTMLURL:     input.HTMLURL,
		ExternalRef: input.ExternalRef,
		TriggeredBy: input.TriggeredBy,
	}).RETURNING(releasesTable.AllColumns)

	query, args := stmt.Sql()

	var model Release
	if err := r.execer.GetContext(
		ctx,
		&model,
		query,
		args...,
	); err != nil {
		return nil, misc.WrapError(err, errs.NewDatabaseError("error while creating release", err.Error()))
	}

	release = model.ToEntity()
	if release == nil {
		return nil, errs.NewInternalServerError("failed to convert release model to entity", nil)
	}

	return release, nil
}
