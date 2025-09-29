package models

// Function represents an executable unit that can be triggered by the
// execution service. Each function has a name and an execution handler that
// defines its behavior.
type Function struct {
	// Name is the unique identifier of the function. It must match the
	// `FunctionName` field of a Trigger for the function to be invoked.
	Name string `json:"name"`
	// Exec defines the logic to be executed when the function is triggered. It
	// receives a pointer to the Trigger that caused the invocation and returns
	// an error if execution fails. As a side effect, it manages the
	// `Execution` model lifecycle (creation, status updates, etc.).
	Exec func(trigger *Trigger) error
}
