package models

// ExecFunction represents an executable unit that can be triggered by the
// execution service.
type ExecFunction func(trigger *Trigger) error
