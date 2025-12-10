package projectrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type scaffoldRequestRepositoryImpl struct {
	execer db.SqlExecer
}

func NewScaffoldRequestRepository(execer db.SqlExecer) repository.ScaffoldRequestRepository {
	return &scaffoldRequestRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *scaffoldRequestRepositoryImpl) WithTx(tx db.SqlExecer) repository.ScaffoldRequestRepository {
	return &scaffoldRequestRepositoryImpl{execer: tx}
}
