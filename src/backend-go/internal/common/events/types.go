package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/platepilot/backend/internal/common/dto"
)

// Event is the base interface for all domain events
type Event interface {
	EventID() uuid.UUID
	EventType() string
	OccurredAt() time.Time
	AggregateID() uuid.UUID
}

// BaseEvent provides common event fields
type BaseEvent struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	OccurredOn  time.Time `json:"occurredOn"`
	AggregateId uuid.UUID `json:"aggregateId"`
}

func (e BaseEvent) EventID() uuid.UUID     { return e.ID }
func (e BaseEvent) EventType() string      { return e.Type }
func (e BaseEvent) OccurredAt() time.Time  { return e.OccurredOn }
func (e BaseEvent) AggregateID() uuid.UUID { return e.AggregateId }

// RecipeCreatedEvent is published when a new recipe is created
type RecipeCreatedEvent struct {
	BaseEvent
	Recipe dto.RecipeDTO `json:"recipe"`
}

// NewRecipeCreatedEvent creates a new RecipeCreatedEvent
func NewRecipeCreatedEvent(recipe dto.RecipeDTO) RecipeCreatedEvent {
	return RecipeCreatedEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New(),
			Type:        "RecipeCreatedEvent",
			OccurredOn:  time.Now().UTC(),
			AggregateId: recipe.ID,
		},
		Recipe: recipe,
	}
}

// RecipeUpdatedEvent is published when a recipe is updated
type RecipeUpdatedEvent struct {
	BaseEvent
}

// NewRecipeUpdatedEvent creates a new RecipeUpdatedEvent
func NewRecipeUpdatedEvent(recipeID uuid.UUID) RecipeUpdatedEvent {
	return RecipeUpdatedEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New(),
			Type:        "RecipeUpdatedEvent",
			OccurredOn:  time.Now().UTC(),
			AggregateId: recipeID,
		},
	}
}
