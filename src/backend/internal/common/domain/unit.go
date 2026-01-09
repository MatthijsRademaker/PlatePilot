package domain

import (
	"time"

	"github.com/google/uuid"
)

// Unit represents a measurement unit for recipe ingredients.
type Unit struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
}
