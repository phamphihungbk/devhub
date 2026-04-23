package approvalrepo

import (
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/db"
)

type approvalRepositoryImpl struct {
	execer db.SqlExecer
}

func NewApprovalRepository(execer db.SqlExecer) repository.ApprovalRepository {
	return &approvalRepositoryImpl{execer: execer}
}
