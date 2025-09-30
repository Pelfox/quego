package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Pelfox/quego/internal/repositories"
	"github.com/Pelfox/quego/models"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

// ErrFunctionNotFound is returned when an attempt is made to process a trigger
// whose target function has not been registered with the `ExecutionService`.
var ErrFunctionNotFound = errors.New("the requested function is not registered")

// ExecutionService provides operations related to `Execution` entities. It
// uses an `ExecutionRepository` for data persistence while serving as the main
// access point for higher layers.
type ExecutionService struct {
	redis     *redis.Client
	repo      *repositories.ExecutionRepository
	functions map[string]models.ExecFunction
	workerSem chan struct{}
}

// NewExecutionService creates and returns a new `ExecutionService` instance
// backed by the provided `ExecutionRepository` and a Redis client.
func NewExecutionService(
	redis *redis.Client,
	repo *repositories.ExecutionRepository,
) *ExecutionService {
	return &ExecutionService{
		redis:     redis,
		repo:      repo,
		functions: make(map[string]models.ExecFunction),
		workerSem: make(chan struct{}, 1),
	}
}

// RegisterFunction adds a new `Function` to the service in order. Registered
// functions can later be invoked or managed by the `ExecutionService`.
func (s *ExecutionService) RegisterFunction(name string, f models.ExecFunction) {
	s.functions[name] = f
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
	_, ok := s.functions[trigger.FunctionName]
	if !ok {
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

	model := models.ExecutionWithTrigger{
		Execution: *payload,
		Trigger:   *trigger,
	}

	data, err := json.Marshal(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trigger: %w", err)
	}

	if err := s.redis.LPush(context.Background(), "quego:queue", data).Err(); err != nil {
		return nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	// go func() {
	// 	if err := s.repo.UpdateStatus(payload.ID, models.ExecutionStatusRunning); err != nil {
	// 		log.Errorf("Failed to update status for job %s: %v", payload.ID, err)
	// 		return
	// 	}
	// 	if err := f(trigger); err != nil {
	// 		log.Errorf("Function execution failed for job %s: %v", payload.ID, err)
	// 		if err := s.repo.UpdateStatus(payload.ID, models.ExecutionStatusFailed); err != nil {
	// 			log.Errorf("Failed to update status for job %s: %v", payload.ID, err)
	// 		}
	// 		return
	// 	}
	// 	if err := s.repo.UpdateStatus(payload.ID, models.ExecutionStatusCompleted); err != nil {
	// 		log.Errorf("Failed to update status for job %s: %v", payload.ID, err)
	// 		return
	// 	}
	// }()

	return payload, nil
}

func (s *ExecutionService) StartWorkers(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := s.redis.BLPop(ctx, 0, "quego:queue").Result()
				if err != nil {
					log.Errorf("Failed to dequeue job: %v", err)
					continue
				}

				var payload models.ExecutionWithTrigger
				if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
					log.Errorf("Failed to unmarshal job payload: %v", err)
					continue
				}

				f, ok := s.functions[payload.Trigger.FunctionName]
				if !ok {
					log.Errorf("Function not found: %s", payload.Trigger.FunctionName)
					continue
				}

				s.workerSem <- struct{}{}
				if err := s.repo.UpdateStatus(payload.Execution.ID, models.ExecutionStatusRunning); err != nil {
					<-s.workerSem
					log.Errorf("Failed to update status for job %s: %v", payload.Execution.ID, err)
					continue
				}

				if err := f(&payload.Trigger); err != nil {
					<-s.workerSem
					log.Errorf("Function execution failed for job %s: %v", payload.Execution.ID, err)
					if err := s.repo.UpdateStatus(payload.Execution.ID, models.ExecutionStatusFailed); err != nil {
						log.Errorf("Failed to update status for job %s: %v", payload.Execution.ID, err)
					}
					continue
				}
				if err := s.repo.UpdateStatus(payload.Execution.ID, models.ExecutionStatusCompleted); err != nil {
					<-s.workerSem
					log.Errorf("Failed to update status for job %s: %v", payload.Execution.ID, err)
					continue
				}
				<-s.workerSem
			}
		}
	}()
}

// GetByID retrieves an `Execution` entity by its unique identifier. It
// delegates the lookup to the underlying repository.
func (s *ExecutionService) GetByID(id uuid.UUID) (*models.Execution, error) {
	return s.repo.GetByID(id)
}

// ListAllTriggers retrieves all `Execution` entities from the underlying repository.
func (s *ExecutionService) ListAllTriggers() ([]*models.ExecutionWithTrigger, error) {
	return s.repo.ListAll()
}

func (s *ExecutionService) RequeueStaled() error {
	staled, err := s.repo.GetStaled()
	if err != nil {
		return err
	}

	for _, exec := range staled {
		data, err := json.Marshal(&models.ExecutionWithTrigger{
			Execution: models.Execution{
				ID:         exec.Execution.ID,
				Status:     exec.Execution.Status,
				TriggerID:  exec.TriggerID,
				StartedAt:  exec.StartedAt,
				FinishedAt: exec.FinishedAt,
			},
			Trigger: models.Trigger{
				TriggerType:  exec.Trigger.TriggerType,
				FunctionName: exec.Trigger.FunctionName,
				Payload:      exec.Trigger.Payload,
				ID:           exec.Trigger.ID,
			},
		})
		if err != nil {
			log.Errorf("Failed to marshal staled execution %s: %v", exec.Execution.ID, err)
			continue
		}

		if err := s.redis.LPush(context.Background(), "quego:queue", data).Err(); err != nil {
			log.Errorf("Failed to re-enqueue staled execution %s: %v", exec.Execution.ID, err)
			continue
		}

		if err := s.repo.UpdateStatus(exec.Execution.ID, models.ExecutionStatusPending); err != nil {
			log.Errorf("Failed to update status for staled execution %s: %v", exec.Execution.ID, err)
			continue
		}
	}

	return nil
}
