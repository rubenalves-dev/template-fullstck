package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/domain"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/httputil"
)

type pgxRepo struct {
	pool *pgxpool.Pool
}

func NewPgxRepository(pool *pgxpool.Pool) domain.Repository {
	return &pgxRepo{pool: pool}
}

func (r *pgxRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, full_name FROM users WHERE email = $1`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httputil.ErrNotFound
		}
		return nil, fmt.Errorf("auth repo get user by email: %w", err)
	}

	return &user, nil
}

func (r *pgxRepo) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, password_hash, full_name) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.FullName)
	if err != nil {
		return fmt.Errorf("auth repo create user: %w", err)
	}
	return nil
}
