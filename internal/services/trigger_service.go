package services

import (
	"github.com/Pelfox/quego/internal/repositories"
	"github.com/Pelfox/quego/models"
	"github.com/google/uuid"
)

// TriggerService provides operations related to `Trigger` entities. It
// uses an `TriggerRepository` for data persistence while serving as the main
// access point for higher layers.
type TriggerService struct {
	repo *repositories.TriggerRepository
}

// NewTriggerService creates and returns a new `TriggerService` instance backed
// by the provided `TriggerRepository`.
func NewTriggerService(repo *repositories.TriggerRepository) *TriggerService {
	return &TriggerService{repo: repo}
}

// Create persists a new `Trigger` in the underlying repository. The provided
// `Trigger` may not have an ID assigned yet; in such cases, the service is
// responsible for generating and assigning a unique ID before storing it.
func (s *TriggerService) Create(trigger *models.Trigger) error {
	triggerID := uuid.New()
	trigger.ID = &triggerID
	return s.repo.Create(trigger)
}
