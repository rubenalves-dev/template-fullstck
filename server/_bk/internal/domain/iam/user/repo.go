package user

import "context"

type CreateUserParams struct {
	Email        string
	PasswordHash string
	FullName     string
}

type CreateOrganizationParams struct {
	Name string
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUserWithOrganization(ctx context.Context, user CreateUserParams, organization CreateOrganizationParams) error
}
