package organizationuser

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationUser struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           UserRole
	CreatedAt      time.Time
}

func New(organizationID, userID uuid.UUID, role UserRole, acceptedAt, createdAt *time.Time) (*OrganizationUser, error) {
	if !role.IsValid() {
		return nil, ErrInvalidUserRole
	}

	createdAtValue := time.Now()
	if createdAt != nil {
		createdAtValue = *createdAt
	}

	return &OrganizationUser{
		OrganizationID: organizationID,
		UserID:         userID,
		Role:           role,
		CreatedAt:      createdAtValue,
	}, nil
}
