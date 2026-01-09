package domain

import (
	"time"

	"github.com/google/uuid"
)

// Ingredient represents an ingredient that can be used in recipes
type Ingredient struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Allergies   []Allergy
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
