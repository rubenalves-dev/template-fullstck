package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/common"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/iam/user"
	"github.com/rubenalves-dev/template-fullstack/server/internal/infrastructure/persistence/postgres/sqlc"
)

type UserRepository struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
		q:    sqlc.New(pool),
	}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	u := &user.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName,
		CreatedAt:    timeOrZero(nil), // CreatedAt is not returned by GetUserByEmail query, assumes standard fields
		UpdatedAt:    timeOrZero(nil),
		ActivatedAt:  timeOrZero(&row.ActivatedAt),
		ArchivedAt:   timeOrZero(&row.ArchivedAt),
	}

	if row.OrganizationID.Valid {
		u.OrganizationID = row.OrganizationID.Bytes
	}
	if row.UserRole.Valid {
		u.Role = string(row.UserRole.UserRole)
	}

	return u, nil
}

func (r *UserRepository) CreateUserWithOrganization(ctx context.Context, u user.CreateUserParams, o user.CreateOrganizationParams) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	// 1. Create User
	newUser, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 2. Create Organization
	newOrg, err := qtx.CreateOrganization(ctx, sqlc.CreateOrganizationParams{
		Name: o.Name,
		Slug: o.Name, // TODO: generate slug properly
	})
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	// 3. Link User to Organization
	err = qtx.LinkUserToOrganization(ctx, sqlc.LinkUserToOrganizationParams{
		OrganizationID: newOrg.ID,
		UserID:         newUser.ID,
		Role:           sqlc.UserRoleADMIN,
	})
	if err != nil {
		return fmt.Errorf("failed to link user to organization: %w", err)
	}

	return tx.Commit(ctx)
}

func timeOrZero(ts *pgtype.Timestamptz) (t time.Time) {
	if ts != nil && ts.Valid {
		return ts.Time
	}
	return
}

// Ensure interface implementation
var _ user.Repository = (*UserRepository)(nil)
