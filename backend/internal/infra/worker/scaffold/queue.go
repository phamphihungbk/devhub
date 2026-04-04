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
