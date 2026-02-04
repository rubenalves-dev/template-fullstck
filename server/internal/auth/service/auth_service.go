package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/httputil"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo      auth.Repository
	nc        *nats.Conn
	jwtSecret string
}

func NewAuthService(repository auth.Repository, nc *nats.Conn, jwtSecret string) auth.Service {
	return &authService{
		repo:      repository,
		nc:        nc,
		jwtSecret: jwtSecret,
	}
}

func (a authService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, httputil.ErrNotFound) {
			return "", httputil.ErrUnauthorized // Don't reveal user existence
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", httputil.ErrUnauthorized
	}

	claims := auth.UserClaims{
		UserID: u.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   u.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

func (a authService) Register(ctx context.Context, user auth.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return a.repo.CreateUser(ctx, &user)
}
