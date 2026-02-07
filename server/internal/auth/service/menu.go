package service

import (
	"sort"

	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/domain"
)

func buildMenuTree(defs []domain.MenuDefinition, userPerms map[string]bool) []domain.MenuNode {
	childrenByParent := make(map[string][]domain.MenuDefinition)
	rootKey := ""
	for _, d := range defs {
		parent := d.ParentID
		childrenByParent[parent] = append(childrenByParent[parent], d)
	}

	var build func(parentID string) []domain.MenuNode
	build = func(parentID string) []domain.MenuNode {
		children := childrenByParent[parentID]
		sort.Slice(children, func(i, j int) bool {
			if children[i].Order == children[j].Order {
				return children[i].Label < children[j].Label
			}
			return children[i].Order < children[j].Order
		})

		var nodes []domain.MenuNode
		for _, d := range children {
			if !d.Visible {
				continue
			}
			if len(d.Permissions) > 0 && !hasAnyPermission(d.Permissions, userPerms) {
				continue
			}

			sub := build(d.ID)
			// If it's a leaf node (no children) and has no path, it's probably just a container
			// that should only be visible if it has visible children.
			// BUT, if it has a path, it's a clickable menu item and should be shown even without children.
			if len(sub) == 0 && d.Path == "" {
				continue
			}

			node := domain.MenuNode{
				Label:    d.Label,
				Path:     d.Path,
				Icon:     d.Icon,
				Children: sub,
			}
			nodes = append(nodes, node)
		}
		return nodes
	}

	return build(rootKey)
}

func hasAnyPermission(perms []string, userPerms map[string]bool) bool {
	for _, p := range perms {
		if userPerms[p] {
			return true
		}
	}
	return false
}
