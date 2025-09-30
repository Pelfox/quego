package models

import (
	"github.com/google/uuid"
)

// TriggerType represents the category of a trigger. A trigger defines the
// condition or mechanism that initiates the execution of a function.
type TriggerType string

const (
	// TriggerTypeCron represents a CRON-based trigger. This type is used when
	// a function should be executed on a recurring schedule defined with CRON
	// syntax.
	TriggerTypeCron TriggerType = "CRON"
	// TriggerTypeEvent represents an event-based trigger. This type is used
	// when a function should be executed in response to an external event.
	TriggerTypeEvent TriggerType = "EVENT"
)

// Trigger describes a request to execute a function. It contains the trigger
// type, the name of the function to be invoked, and payload data to be passed
// into the function.
type Trigger struct {
	// ID is the unique identifier for the trigger. It will be `nil`
	// immediately after a trigger is received. Before function execution, it
	// will be assigned by the persistence layer.
	ID *uuid.UUID `db:"id" json:"id,omitempty"`
	// TriggerType specifies how the trigger was created and from where it was
	// received.
	TriggerType TriggerType `db:"trigger_type" json:"trigger_type"`
	// FunctionName is the name of the function to be executed. This must match
	// the name of a function that has been registered with the execution
	// service.
	FunctionName string `db:"function_name" json:"function_name"`
	// Payload contains the input data for the function execution. It must
	// match the defined input schema of the target function.
	Payload string `db:"payload" json:"-"`
}
