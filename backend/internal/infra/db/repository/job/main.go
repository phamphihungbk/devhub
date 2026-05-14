package jobrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type jobRepositoryImpl struct {
	execer db.SqlExecer
}

func NewJobRepository(execer db.SqlExecer) repository.JobRepository {
	return &jobRepositoryImpl{execer: execer}
}

func (r *jobRepositoryImpl) WithTx(tx db.SqlExecer) repository.JobRepository {
	return &jobRepositoryImpl{execer: tx}
}
