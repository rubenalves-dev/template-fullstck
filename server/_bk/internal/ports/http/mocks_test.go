package http

import (
	"context"
	"github.com/rubenalves-dev/template-fullstack/server/internal/applications/auth"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, params auth.RegisterParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}
