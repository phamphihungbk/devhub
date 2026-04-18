package service

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type serviceRepositoryImpl struct {
	execer db.SqlExecer
}

func NewServiceRepository(execer db.SqlExecer) repository.ServiceRepository {
	return &serviceRepositoryImpl{execer: execer}
}
