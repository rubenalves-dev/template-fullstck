package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/common"
	iam_auth "github.com/rubenalves-dev/template-fullstack/server/internal/domain/iam/auth"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/iam/user"
	"golang.org/x/crypto/bcrypt"
)

type RegisterParams struct {
	Email            string
	Password         string
	FullName         string
	OrganizationName string
}

type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, params RegisterParams) error
}

type service struct {
	userRepo  user.Repository
	jwtSecret string
}

func NewService(userRepo user.Repository, jwtSecret string) Service {
	return &service{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return "", common.ErrUnauthorized // Don't reveal user existence
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", common.ErrUnauthorized
	}

	// Generate Token
	claims := iam_auth.UserClaims{
		UserID:         u.ID.String(),
		OrganizationID: u.OrganizationID.String(),
		Role:           u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   u.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *service) Register(ctx context.Context, params RegisterParams) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.CreateUserWithOrganization(ctx,
		user.CreateUserParams{
			Email:        params.Email,
			PasswordHash: string(hashedPassword),
			FullName:     params.FullName,
		},
		user.CreateOrganizationParams{
			Name: params.OrganizationName,
		},
	)
}
