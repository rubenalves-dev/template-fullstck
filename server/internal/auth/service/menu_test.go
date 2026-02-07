package service

import (
	"reflect"
	"testing"

	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/domain"
)

func TestBuildMenuTreeFiltersAndOrders(t *testing.T) {
	defs := []domain.MenuDefinition{
		{ID: "dashboard", Label: "Dashboard", Path: "/dashboard", Order: 0, Visible: true}, // No permissions
		{ID: "root", Label: "Root", Order: 10, Visible: true},
		{ID: "a", Label: "A", ParentID: "root", Path: "/a", Order: 20, Visible: true},
		{ID: "b", Label: "B", ParentID: "root", Path: "/b", Order: 10, Permissions: []string{"p.read"}, Visible: true},
		{ID: "c", Label: "C", ParentID: "root", Path: "/c", Order: 30, Permissions: []string{"p.hidden"}, Visible: true},
	}

	userPerms := map[string]bool{"p.read": true}
	menu := buildMenuTree(defs, userPerms)

	expected := []domain.MenuNode{
		{
			Label: "Dashboard",
			Path:  "/dashboard",
		},
		{
			Label: "Root",
			Children: []domain.MenuNode{
				{Label: "B", Path: "/b"},
				{Label: "A", Path: "/a"},
			},
		},
	}

	if !reflect.DeepEqual(menu, expected) {
		t.Fatalf("unexpected menu: %#v", menu)
	}
}
