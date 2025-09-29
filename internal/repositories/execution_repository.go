package repositories

import (
	"time"

	"github.com/Pelfox/quego/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ExecutionRepository handles database operations for `Execution` entities.
type ExecutionRepository struct {
	db *sqlx.DB
}

// NewExecutionRepository creates a new `ExecutionRepository` backed by the
// given `sqlx.DB` instance.
func NewExecutionRepository(db *sqlx.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

// Create inserts a new `Execution` record into the database. The
// provided `Execution` struct must include values for `id`, `status`,
// and `trigger_id`.
func (r *ExecutionRepository) Create(data *models.Execution) error {
	query := "INSERT INTO executions (id, status, trigger_id) VALUES (:id, :status, :trigger_id)"
	_, err := r.db.NamedExec(query, data)
	return err
}

// UpdateStatus updates the status of an `Execution` model in the database. In
// addition to the status field, it conditionally updates timestamp fields
// depending on the new status:
// - `ExecutionStatusRunning`: updates `started_at`.
// - `ExecutionStatusCompleted` or `ExecutionStatusFailed`: updates `finished_at`.
func (r *ExecutionRepository) UpdateStatus(id uuid.UUID, newStatus models.ExecutionStatus) error {
	var (
		query string
		args  []any
	)

	switch newStatus {
	case models.ExecutionStatusRunning:
		query = "UPDATE executions SET status = ?, started_at = ? WHERE id = ?"
		args = []any{newStatus, time.Now(), id}
	case models.ExecutionStatusCompleted, models.ExecutionStatusFailed:
		query = "UPDATE executions SET status = ?, finished_at = ? WHERE id = ?"
		args = []any{newStatus, time.Now(), id}
	default:
		query = "UPDATE executions SET status = ? WHERE id = ?"
		args = []any{newStatus, id}
	}

	_, err := r.db.Exec(query, args...)
	return err
}

// GetByID retrieves an `Execution` model by its unique identifier.
func (r *ExecutionRepository) GetByID(id uuid.UUID) (*models.Execution, error) {
	var execution models.Execution
	query := "SELECT * FROM executions WHERE id = ?"
	if err := r.db.Get(&execution, query, id); err != nil {
		return nil, err
	}
	return &execution, nil
}
