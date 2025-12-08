package projectrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type projectRepositoryImpl struct {
	execer db.SqlExecer
}

func NewProjectRepository(execer db.SqlExecer) repository.ProjectRepository {
	return &projectRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *projectRepositoryImpl) WithTx(tx db.SqlExecer) repository.ProjectRepository {
	return &projectRepositoryImpl{execer: tx}
}
