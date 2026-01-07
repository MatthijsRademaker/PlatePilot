package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated user.
type User struct {
	ID          uuid.UUID
	Email       string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
