package domain

const (
	PermissionPageRead   = "cms.page.read"
	PermissionPageWrite  = "cms.page.write"
	PermissionPageDelete = "cms.page.delete"
)

func GetAvailablePermission() []string {
	return []string{PermissionPageRead, PermissionPageWrite, PermissionPageDelete}
}
