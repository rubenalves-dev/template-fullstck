package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func (r *pgxRepo) UpsertPermissions(ctx context.Context, permissions []domain.Permission) error {
	if len(permissions) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, p := range permissions {
		batch.Queue(`
			INSERT INTO permissions (id, module, description) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (id) DO UPDATE SET module = EXCLUDED.module, description = EXCLUDED.description`,
			p.ID, p.Module, p.Description,
		)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("auth repo upsert permissions: %w", err)
	}
	return nil
}

func (r *pgxRepo) CreateRole(ctx context.Context, name string) (*domain.Role, error) {
	query := `INSERT INTO roles (name) VALUES ($1) RETURNING id, name`
	var role domain.Role
	err := r.pool.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, fmt.Errorf("auth repo create role: %w", err)
	}
	return &role, nil
}

func (r *pgxRepo) GetRoles(ctx context.Context) ([]domain.Role, error) {
	query := `SELECT id, name FROM roles ORDER BY id`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("auth repo get roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *pgxRepo) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("auth repo assign role: %w", err)
	}
	return nil
}

func (r *pgxRepo) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
		SELECT DISTINCT p.id
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("auth repo get user permissions: %w", err)
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *pgxRepo) AddPermissionToRole(ctx context.Context, roleID int, permissionID string) error {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("auth repo add permission to role: %w", err)
	}
	return nil
}

func (r *pgxRepo) UpsertMenuDefinitions(ctx context.Context, defs []domain.MenuDefinition) error {
	if len(defs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, d := range defs {
		permissions := d.Permissions
		if permissions == nil {
			permissions = []string{}
		}
		batch.Queue(`
			INSERT INTO menu_definitions
				(id, domain, label, path, icon, order_index, parent_id, permissions, visible, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
			ON CONFLICT (id) DO UPDATE SET
				domain = EXCLUDED.domain,
				label = EXCLUDED.label,
				path = EXCLUDED.path,
				icon = EXCLUDED.icon,
				order_index = EXCLUDED.order_index,
				parent_id = EXCLUDED.parent_id,
				permissions = EXCLUDED.permissions,
				visible = EXCLUDED.visible,
				updated_at = now()
		`, d.ID, d.Domain, d.Label, d.Path, d.Icon, d.Order, nullableString(d.ParentID), permissions, d.Visible)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("auth repo upsert menu definitions: %w", err)
	}
	return nil
}

func (r *pgxRepo) GetMenuDefinitions(ctx context.Context) ([]domain.MenuDefinition, error) {
	query := `
		SELECT id, domain, label, path, icon, order_index, parent_id, permissions, visible
		FROM menu_definitions
		ORDER BY domain, order_index, id
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("auth repo get menu definitions: %w", err)
	}
	defer rows.Close()

	var defs []domain.MenuDefinition
	for rows.Next() {
		var d domain.MenuDefinition
		var parentID *string
		if err := rows.Scan(&d.ID, &d.Domain, &d.Label, &d.Path, &d.Icon, &d.Order, &parentID, &d.Permissions, &d.Visible); err != nil {
			return nil, err
		}
		if parentID != nil {
			d.ParentID = *parentID
		}
		defs = append(defs, d)
	}
	return defs, nil
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
