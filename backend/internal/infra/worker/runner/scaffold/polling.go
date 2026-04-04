package scaffold_runner

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TODO: consider reusing repository or generated DB models for this job shape.
type ScaffoldJob struct {
	ID          uuid.UUID
	Template    string
	ProjectID   uuid.UUID
	Environment string
	Variables   string
}

func (j ScaffoldJob) GetID() uuid.UUID {
	return j.ID
}

type ScaffoldQueueSourceAdapter struct {
	db *sqlx.DB
}

func NewScaffoldQueueSourceAdapter(db *sqlx.DB) *ScaffoldQueueSourceAdapter {
	return &ScaffoldQueueSourceAdapter{db: db}
}

func (a *ScaffoldQueueSourceAdapter) Dequeue(ctx context.Context) (*ScaffoldJob, error) {
	var row struct {
		ID          uuid.UUID `db:"id"`
		Template    string    `db:"template"`
		ProjectID   uuid.UUID `db:"project_id"`
		Environment string    `db:"environment"`
		Variables   string    `db:"variables"`
	}

	// TODO: use repository instead of raw SQL here.
	const query = `
SELECT id, template, project_id, environment, variables
FROM scaffold_requests
WHERE status = 'pending'
ORDER BY created_at ASC, id ASC
LIMIT 1`

	if err := a.db.GetContext(ctx, &row, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("dequeue scaffold request: %w", err)
	}

	return &ScaffoldJob{
		ID:          row.ID,
		Template:    row.Template,
		ProjectID:   row.ProjectID,
		Environment: row.Environment,
		Variables:   row.Variables,
	}, nil
}

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
