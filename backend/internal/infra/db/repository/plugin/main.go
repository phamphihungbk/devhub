package pluginrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type pluginRepositoryImpl struct {
	execer db.SqlExecer
}

func NewPluginRepository(execer db.SqlExecer) repository.PluginRepository {
	return &pluginRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *pluginRepositoryImpl) WithTx(tx db.SqlExecer) repository.PluginRepository {
	return &pluginRepositoryImpl{execer: tx}
}
