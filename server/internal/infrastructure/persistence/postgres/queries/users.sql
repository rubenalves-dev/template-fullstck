-- name: GetUserByEmail :one
SELECT 
    u.id, 
    u.email, 
    u.password_hash, 
    u.full_name, 
    u.activated_at, 
    u.archived_at,
    ou.organization_id,
    ou.role as user_role
FROM users u
LEFT JOIN organization_users ou ON u.id = ou.user_id
WHERE u.email = $1 AND u.archived_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, full_name
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CreateOrganization :one
INSERT INTO organizations (
    name, slug
) VALUES (
    $1, $2
) RETURNING *;

-- name: LinkUserToOrganization :exec
INSERT INTO organization_users (
    organization_id, user_id, role
) VALUES (
    $1, $2, $3
);
