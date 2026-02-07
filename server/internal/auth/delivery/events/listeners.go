package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/domain"
	"github.com/rubenalves-dev/template-fullstack/server/internal/platform/menu"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/events"
)

type eventHandler struct {
	svc domain.Service
}

func RegisterListeners(nc *nats.Conn, svc domain.Service) {
	h := &eventHandler{svc: svc}

	// Example
	_, err := nc.Subscribe("ecommerce.order.completed", h.handleOrderCompleted)
	if err != nil {
		log.Printf("Failed to subscribe to ecommerce.order.completed: %v", err)
	}

	_, err = nc.Subscribe(events.SystemPermissionsRegister, h.handlePermissionsRegister)
	if err != nil {
		log.Printf("Failed to subscribe to %s: %v", events.SystemPermissionsRegister, err)
	}

	_, err = nc.Subscribe(events.SystemMenusRegister, h.handleMenusRegister)
	if err != nil {
		log.Printf("Failed to subscribe to %s: %v", events.SystemMenusRegister, err)
	}
}

func (h *eventHandler) handlePermissionsRegister(m *nats.Msg) {
	var payload events.SystemPermissionsRegisteredData
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		log.Printf("Failed to unmarshal permission event: %v", err)
		return
	}

	if err := h.svc.RegisterModulePermissions(context.Background(), payload.Module, payload.Permissions); err != nil {
		log.Printf("Failed to register permissions for module %s: %v", payload.Module, err)
	}
}

func (h *eventHandler) handleMenusRegister(m *nats.Msg) {
	var payload events.SystemMenusRegisteredData
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		log.Printf("Failed to unmarshal menus event: %v", err)
		return
	}

	defs := flattenMenuDefinitions(payload.Domain, payload.Menu, "")
	if err := h.svc.RegisterModuleMenus(context.Background(), payload.Domain, defs); err != nil {
		log.Printf("Failed to register menus for domain %s: %v", payload.Domain, err)
	}
}

func flattenMenuDefinitions(domainName string, nodes []menu.MenuDefinition, parentID string) []domain.MenuDefinition {
	var defs []domain.MenuDefinition
	for _, n := range nodes {
		def := n
		def.Domain = domainName
		def.ParentID = parentID
		// Recursively flatten children
		children := n.Children
		def.Children = nil // Clear children for the flattened definition
		defs = append(defs, def)

		if len(children) > 0 {
			defs = append(defs, flattenMenuDefinitions(domainName, children, n.ID)...)
		}
	}
	return defs
}

func (h *eventHandler) handleOrderCompleted(m *nats.Msg) {
	var payload struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(m.Data, &payload); err != nil {
		log.Printf("Failed to unmarshal event payload: %v", err)
	}

	log.Printf("Auth Module: User %s has completed an order.", payload.UserID)
}
