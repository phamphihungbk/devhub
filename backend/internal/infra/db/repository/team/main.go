package teamrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type teamRepositoryImpl struct {
	execer db.SqlExecer
}

func NewTeamRepository(execer db.SqlExecer) repository.TeamRepository {
	return &teamRepositoryImpl{execer: execer}
}
