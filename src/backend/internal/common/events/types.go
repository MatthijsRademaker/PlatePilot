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
	ID               uuid.UUID `json:"id"`
	Type             string    `json:"type"`
	OccurredOn        time.Time `json:"occurredOn"`
	AggregateId       uuid.UUID `json:"aggregateId"`
	SchemaVersion     int       `json:"schemaVersion"`
	AggregateVersion  int       `json:"aggregateVersion"`
}

func (e BaseEvent) EventID() uuid.UUID     { return e.ID }
func (e BaseEvent) EventType() string      { return e.Type }
func (e BaseEvent) OccurredAt() time.Time  { return e.OccurredOn }
func (e BaseEvent) AggregateID() uuid.UUID { return e.AggregateId }

// RecipeCreatedEvent is published when a new recipe is created
type RecipeUpsertedEvent struct {
	BaseEvent
	Recipe dto.RecipeDTO `json:"recipe"`
}

// NewRecipeUpsertedEvent creates a new RecipeUpsertedEvent.
func NewRecipeUpsertedEvent(recipe dto.RecipeDTO) RecipeUpsertedEvent {
	return RecipeUpsertedEvent{
		BaseEvent: BaseEvent{
			ID:              uuid.New(),
			Type:            "RecipeUpsertedEvent",
			OccurredOn:       time.Now().UTC(),
			AggregateId:      recipe.ID,
			SchemaVersion:    1,
			AggregateVersion: 0,
		},
		Recipe: recipe,
	}
}

// RecipeDeletedEvent is published when a recipe is deleted.
type RecipeDeletedEvent struct {
	BaseEvent
	UserID    uuid.UUID `json:"userId"`
	DeletedAt time.Time `json:"deletedAt"`
}

// NewRecipeDeletedEvent creates a new RecipeDeletedEvent.
func NewRecipeDeletedEvent(recipeID, userID uuid.UUID) RecipeDeletedEvent {
	return RecipeDeletedEvent{
		BaseEvent: BaseEvent{
			ID:              uuid.New(),
			Type:            "RecipeDeletedEvent",
			OccurredOn:       time.Now().UTC(),
			AggregateId:      recipeID,
			SchemaVersion:    1,
			AggregateVersion: 0,
		},
		UserID:    userID,
		DeletedAt: time.Now().UTC(),
	}
}
