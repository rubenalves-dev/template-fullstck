package auth

import "github.com/rubenalves-dev/template-fullstack/server/internal/auth/domain"

var MenuDefinitions = []domain.MenuDefinition{
	{
		ID:      "core:dashboard",
		Label:   "Dashboard",
		Path:    "/dashboard",
		Icon:    "dashboard",
		Order:   0,
		Visible: true,
	},
	{
		ID:          "auth:system",
		Label:       "System",
		Icon:        "settings",
		Order:       90,
		Permissions: []string{domain.PermissionRoleRead},
		Visible:     true,
	},
	{
		ID:          "auth:roles",
		Label:       "Roles",
		Path:        "/system/roles",
		Order:       10,
		ParentID:    "auth:system",
		Permissions: []string{domain.PermissionRoleRead},
		Visible:     true,
	},
}
