package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ActivatedAt time.Time
	ArchivedAt  time.Time
}
