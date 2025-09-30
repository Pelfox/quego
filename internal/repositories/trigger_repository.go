package repositories

import (
	"github.com/Pelfox/quego/models"
	"github.com/jmoiron/sqlx"
)

// TriggerRepository handles database operations for `Trigger` entities.
type TriggerRepository struct {
	db *sqlx.DB
}

// NewTriggerRepository creates a new `TriggerRepository` backed by the
// given `sqlx.DB` instance.
func NewTriggerRepository(db *sqlx.DB) *TriggerRepository {
	return &TriggerRepository{db: db}
}

// Create inserts a new `Trigger` record into the database.
func (r *TriggerRepository) Create(data *models.Trigger) error {
	query := "INSERT INTO triggers (id, trigger_type, function_name, payload) VALUES (:id, :trigger_type, :function_name, :payload)"
	_, err := r.db.NamedExec(query, data)
	return err
}
