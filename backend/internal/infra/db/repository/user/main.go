package userrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type userRepositoryImpl struct {
	execer db.SqlExecer
}

func NewUserRepository(execer db.SqlExecer) repository.UserRepository {
	return &userRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *userRepositoryImpl) WithTx(tx db.SqlExecer) repository.UserRepository {
	return &userRepositoryImpl{execer: tx}
}
