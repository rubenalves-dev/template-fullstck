package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents an entity with personal and account-related data managed within the system.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ActivatedAt *time.Time
	ArchivedAt  *time.Time
}

type Permission struct {
	ID          string    `json:"id"`
	Module      string    `json:"module"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
