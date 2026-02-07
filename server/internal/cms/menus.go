package cms

import (
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/domain"
	"github.com/rubenalves-dev/template-fullstack/server/internal/platform/menu"
)

var MenuDefinition = menu.MenuDefinition{
	ID:          "cms:root",
	Label:       "CMS",
	Icon:        "article",
	Order:       20,
	Visible:     true,
	Permissions: []string{domain.PermissionPageRead},
	Children: []menu.MenuDefinition{
		{
			ID:          "cms:pages",
			Label:       "Pages",
			Path:        "/cms/pages",
			Order:       10,
			Visible:     true,
			Permissions: []string{domain.PermissionPageRead},
		},
		{
			ID:          "cms:media",
			Label:       "Media",
			Path:        "/cms/media",
			Order:       20,
			Visible:     true,
			Permissions: []string{domain.PermissionPageRead},
		},
	},
}
