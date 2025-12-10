package userrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type deploymentRepositoryImpl struct {
	execer db.SqlExecer
}

func NewDeploymentRepository(execer db.SqlExecer) repository.DeploymentRepository {
	return &deploymentRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *deploymentRepositoryImpl) WithTx(tx db.SqlExecer) repository.DeploymentRepository {
	return &deploymentRepositoryImpl{execer: tx}
}
