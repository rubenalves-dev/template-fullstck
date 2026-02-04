package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/common"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/iam/user"
	"github.com/rubenalves-dev/template-fullstack/server/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Login(t *testing.T) {
	jwtSecret := "test-secret"
	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	// Generate a valid hash for the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	validUser := &user.User{
		ID:             uuid.New(),
		Email:          email,
		PasswordHash:   string(hashedPassword),
		OrganizationID: uuid.New(),
		Role:           "ADMIN",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		svc := NewService(mockRepo, jwtSecret)

		mockRepo.On("GetUserByEmail", ctx, email).Return(validUser, nil)

		token, err := svc.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		svc := NewService(mockRepo, jwtSecret)

		mockRepo.On("GetUserByEmail", ctx, email).Return(nil, common.ErrNotFound)

		token, err := svc.Login(ctx, email, password)

		assert.ErrorIs(t, err, common.ErrUnauthorized) // Should mask as Unauthorized
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		svc := NewService(mockRepo, jwtSecret)

		mockRepo.On("GetUserByEmail", ctx, email).Return(validUser, nil)

		token, err := svc.Login(ctx, email, "wrongpassword")

		assert.ErrorIs(t, err, common.ErrUnauthorized)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		svc := NewService(mockRepo, jwtSecret)

		expectedErr := errors.New("db error")
		mockRepo.On("GetUserByEmail", ctx, email).Return(nil, expectedErr)

		token, err := svc.Login(ctx, email, password)

		assert.Equal(t, expectedErr, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_Register(t *testing.T) {
	jwtSecret := "test-secret"
	ctx := context.Background()
	params := RegisterParams{
		Email:            "new@example.com",
		Password:         "password123",
		FullName:         "New User",
		OrganizationName: "New Org",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		svc := NewService(mockRepo, jwtSecret)

		mockRepo.On("CreateUserWithOrganization", ctx, mock.MatchedBy(func(u user.CreateUserParams) bool {
			// Verify password is hashed
			err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(params.Password))
			return u.Email == params.Email && u.FullName == params.FullName && err == nil
		}), mock.MatchedBy(func(o user.CreateOrganizationParams) bool {
			return o.Name == params.OrganizationName
		})).Return(nil)

		err := svc.Register(ctx, params)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		svc := NewService(mockRepo, jwtSecret)

		expectedErr := errors.New("db error")
		mockRepo.On("CreateUserWithOrganization", ctx, mock.Anything, mock.Anything).Return(expectedErr)

		err := svc.Register(ctx, params)

		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}