package domain

import (
	"time"

	"github.com/google/uuid"
)

// Cuisine represents a type of cuisine (e.g., Italian, Mexican, etc.)
type Cuisine struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
}
