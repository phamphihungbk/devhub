package refreshtokenrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type tokenRefreshRepositoryImpl struct {
	execer db.SqlExecer
}

func NewRefreshTokenRepository(execer db.SqlExecer) repository.RefreshTokenRepository {
	return &tokenRefreshRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *tokenRefreshRepositoryImpl) WithTx(tx db.SqlExecer) repository.RefreshTokenRepository {
	return &tokenRefreshRepositoryImpl{execer: tx}
}
