package mocks

import (
	"context"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/iam/user"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *UserRepository) CreateUserWithOrganization(ctx context.Context, u user.CreateUserParams, o user.CreateOrganizationParams) error {
	args := m.Called(ctx, u, o)
	return args.Error(0)
}
