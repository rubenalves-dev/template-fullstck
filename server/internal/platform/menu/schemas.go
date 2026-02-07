package menu

type MenuDefinition struct {
	ID          string           `json:"id"`
	Domain      string           `json:"domain"`
	Label       string           `json:"label"`
	Path        string           `json:"path,omitempty"`
	Icon        string           `json:"icon,omitempty"`
	Order       int              `json:"order,omitempty"`
	ParentID    string           `json:"parent_id,omitempty"`
	Permissions []string         `json:"permissions,omitempty"`
	Visible     bool             `json:"visible"`
	Children    []MenuDefinition `json:"children,omitempty"`
}

type MenuNode struct {
	Label    string     `json:"label"`
	Path     string     `json:"path,omitempty"`
	Icon     string     `json:"icon,omitempty"`
	Children []MenuNode `json:"children,omitempty"`
}
