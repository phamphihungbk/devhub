package scaffold_runner

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ScaffoldStatePersistence struct {
	db *sqlx.DB
}

func NewScaffoldStatePersistence(db *sqlx.DB) *ScaffoldStatePersistence {
	return &ScaffoldStatePersistence{db: db}
}

func (p *ScaffoldStatePersistence) MarkRunning(ctx context.Context, id uuid.UUID) error {
	// TODO: use repository instead of raw SQL here.
	const query = `
UPDATE scaffold_requests
SET status = 'running', updated_at = now()
WHERE id = $1 AND status = 'pending'`

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark scaffold request running: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("scaffold request %s is not pending", id)
	}

	return nil
}

func (p *ScaffoldStatePersistence) MarkCompleted(ctx context.Context, id uuid.UUID, result ScaffoldExecutionResult) error {
	// TODO: use repository instead of raw SQL here.
	const query = `
UPDATE scaffold_requests
SET status = 'completed', result_repo_url = $2, updated_at = now()
WHERE id = $1`

	if _, err := p.db.ExecContext(ctx, query, id, result.RepoURL); err != nil {
		return fmt.Errorf("mark scaffold request completed: %w", err)
	}
	return nil
}

func (p *ScaffoldStatePersistence) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	const query = `
UPDATE scaffold_requests
SET status = 'failed', updated_at = now()
WHERE id = $1`

	if _, err := p.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("mark scaffold request failed: %w", err)
	}
	_ = reason
	return nil
}
