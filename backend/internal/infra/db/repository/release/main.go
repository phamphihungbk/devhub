package releaserepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type releaseRepositoryImpl struct {
	execer db.SqlExecer
}

func NewReleaseRepository(execer db.SqlExecer) repository.ReleaseRepository {
	return &releaseRepositoryImpl{execer: execer}
}
