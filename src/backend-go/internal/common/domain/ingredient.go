package domain

import (
	"time"

	"github.com/google/uuid"
)

// Ingredient represents an ingredient that can be used in recipes
type Ingredient struct {
	ID        uuid.UUID
	Name      string
	Quantity  string
	Allergies []Allergy
	CreatedAt time.Time
}
