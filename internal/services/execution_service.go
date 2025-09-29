package services

import (
	"errors"

	"github.com/Pelfox/quego/internal/repositories"
	"github.com/Pelfox/quego/models"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

// ErrFunctionNotFound is returned when an attempt is made to process a trigger
// whose target function has not been registered with the `ExecutionService`.
var ErrFunctionNotFound = errors.New("the requested function is not registered")

// ExecutionService provides operations related to `Execution` entities. It
// uses an `ExecutionRepository` for data persistence while serving as the main
// access point for higher layers.
type ExecutionService struct {
	repo      *repositories.ExecutionRepository
	functions []*models.Function
}

// NewExecutionService creates and returns a new `ExecutionService` instance
// backed by the provided `ExecutionRepository`.
func NewExecutionService(repo *repositories.ExecutionRepository) *ExecutionService {
	return &ExecutionService{repo: repo, functions: make([]*models.Function, 0)}
}

// RegisterFunction adds a new `Function` to the service in order. Registered
// functions can later be invoked or managed by the `ExecutionService`.
func (s *ExecutionService) RegisterFunction(f *models.Function) {
	s.functions = append(s.functions, f)
}

// Process looks up and executes the function associated with call from the
// given `Trigger` payload.
//
// If a matching function is found, its `Exec` method is invoked in a separate
// goroutine, allowing the execution to run asynchronously. The method then
// returns immediately with no error.
//
// If no function matches the trigger's request name, the method returns the
// `ErrFunctionNotFound` error.
func (s *ExecutionService) Process(trigger *models.Trigger) (*models.Execution, error) {
	var targetFunction *models.Function
	for _, function := range s.functions {
		if function.Name == trigger.FunctionName {
			targetFunction = function
			break
		}
	}

	if targetFunction == nil {
		return nil, ErrFunctionNotFound
	}

	payload := &models.Execution{
		ID:        uuid.New(),
		Status:    models.ExecutionStatusPending,
		TriggerID: *trigger.ID,
	}
	if err := s.repo.Create(payload); err != nil {
		return nil, err
	}

	go func() {
		if err := s.repo.UpdateStatus(payload.ID, models.ExecutionStatusRunning); err != nil {
			log.Errorf("Failed to update status for job %s: %v", payload.ID, err)
			return
		}
		if err := targetFunction.Exec(trigger); err != nil {
			log.Errorf("Function execution failed for job %s: %v", payload.ID, err)
			if err := s.repo.UpdateStatus(payload.ID, models.ExecutionStatusFailed); err != nil {
				log.Errorf("Failed to update status for job %s: %v", payload.ID, err)
			}
			return
		}
		if err := s.repo.UpdateStatus(payload.ID, models.ExecutionStatusCompleted); err != nil {
			log.Errorf("Failed to update status for job %s: %v", payload.ID, err)
			return
		}
	}()

	return payload, nil
}

// GetByID retrieves an `Execution` entity by its unique identifier. It
// delegates the lookup to the underlying repository.
func (s *ExecutionService) GetByID(id uuid.UUID) (*models.Execution, error) {
	return s.repo.GetByID(id)
}
