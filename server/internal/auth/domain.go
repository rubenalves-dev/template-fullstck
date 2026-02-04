package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ContextKey string

const UserClaimsKey ContextKey = "user_claims"

// UserClaims represents the claims of a JWT token issued to a user.
type UserClaims struct {
	UserID string
	jwt.RegisteredClaims
}

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

// Repository defines an interface for managing user data storage and retrieval operations in the system.
type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

// Service defines an interface for managing user authentication and registration operations in the system.
type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user User) error
}
