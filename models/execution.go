package models

import (
	"time"

	"github.com/google/uuid"
)

// ExecutionStatus represents the lifecycle state of an Execution. It indicates
// whether the execution is still queued, actively running, finished
// successfully, or ended with an error.
type ExecutionStatus string

const (
	// ExecutionStatusPending means the execution has been registered but
	// has not started yet. This may happen if the system is still queuing
	// the request behind other executions.
	ExecutionStatusPending ExecutionStatus = "PENDING"
	// ExecutionStatusRunning means the execution is currently in progress.
	ExecutionStatusRunning ExecutionStatus = "RUNNING"
	// ExecutionStatusCompleted means the execution finished successfully
	// without errors.
	ExecutionStatusCompleted ExecutionStatus = "COMPLETED"
	// ExecutionStatusFailed means the execution has finished, but with an
	// error or unexpected termination.
	ExecutionStatusFailed ExecutionStatus = "FAILED"
)

// Execution represents a single invocation attempt of a triggered function.
//
// An Execution record is created once an event trigger has been received and
// the target function is defined (i.e. the system knows how to run it).
// Creation of an Execution does not necessarily mean the function has
// started running yet—it may remain in the Pending state until scheduled.
type Execution struct {
	// ID is the unique identifier of this execution.
	ID uuid.UUID `db:"id" json:"id"`
	// Status is the current lifecycle state of this execution.
	Status ExecutionStatus `db:"status" json:"status"`
	// TriggerID refers to the originating trigger that caused this
	// execution to be created.
	TriggerID uuid.UUID `db:"trigger_id" json:"trigger_id"`

	// StartedAt is the timestamp when the execution actually began
	// running. It is nil if the execution has not started yet.
	StartedAt *time.Time `db:"started_at" json:"started_at,omitempty"`
	// FinishedAt is the timestamp when the execution reached a terminal
	// state (`Completed` or `Failed`). It is nil if the execution is still
	// pending or running.
	FinishedAt *time.Time `db:"finished_at" json:"finished_at,omitempty"`
}

// ExecutionWithTrigger represents an execution along with its associated
// trigger details. It is used for queries that need to return both execution
// and trigger information together.
type ExecutionWithTrigger struct {
	Execution
	// Trigger refers to the originating trigger that caused this execution to
	// be created.
	Trigger Trigger `db:"trigger" json:"trigger"`
}

type StaleExecution struct {
	Payload      map[string]any
	FunctionName string
	ExecutionID  uuid.UUID
}
