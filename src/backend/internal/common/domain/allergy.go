package domain

import (
	"time"

	"github.com/google/uuid"
)

// Allergy represents an allergen that may be present in ingredients
type Allergy struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}
