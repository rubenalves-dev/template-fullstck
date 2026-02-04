package events

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth"
)

type eventHandler struct {
	svc auth.Service
}

func RegisterListeners(nc *nats.Conn, svc auth.Service) {
	h := &eventHandler{svc: svc}

	// Example
	_, err := nc.Subscribe("ecommerce.order.completed", h.handleOrderCompleted)
	if err != nil {
		return
	}
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
