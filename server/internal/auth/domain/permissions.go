package domain

const (
	PermissionRoleRead   = "auth.role.read"
	PermissionRoleWrite  = "auth.role.write"
	PermissionRoleDelete = "auth.role.delete"
	PermissionUserRead   = "auth.user.read"
	PermissionUserWrite  = "auth.user.write"
)

func GetAvailablePermissions() []string {
	return []string{
		PermissionRoleRead,
		PermissionRoleWrite,
		PermissionRoleDelete,
		PermissionUserRead,
		PermissionUserWrite,
	}
}
