package environment

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type environmentRepositoryImpl struct {
	execer db.SqlExecer
}

func NewEnvironmentRepository(execer db.SqlExecer) repository.EnvironmentRepository {
	return &environmentRepositoryImpl{execer: execer}
}
