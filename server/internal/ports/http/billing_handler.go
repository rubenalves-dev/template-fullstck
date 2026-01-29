package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/beheer/internal/billing"
	"github.com/rubenalves-dev/beheer/internal/db"
	"github.com/rubenalves-dev/beheer/internal/domain"
)

type BillingHandler struct {
	svc billing.Service
}

func NewBillingHandler(svc billing.Service) *BillingHandler {
	return &BillingHandler{svc: svc}
}

func (h *BillingHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	var req struct {
		Invoice db.CreateInvoiceParams       `json:"invoice"`
		Items   []db.CreateInvoiceItemParams `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}
	req.Invoice.OrganizationID = orgID

	invoice, items, err := h.svc.CreateInvoice(r.Context(), req.Invoice, req.Items)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, map[string]interface{}{
		"invoice": invoice,
		"items":   items,
	})
}

func (h *BillingHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid invoice ID")
		return
	}

	invoice, items, err := h.svc.GetInvoice(r.Context(), id, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, map[string]interface{}{
		"invoice": invoice,
		"items":   items,
	})
}

func (h *BillingHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr != "" {
		memberID, err := uuid.Parse(memberIDStr)
		if err != nil {
			domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid member ID")
			return
		}
		invoices, err := h.svc.ListInvoicesByMember(r.Context(), memberID, orgID)
		if err != nil {
			status, code := domain.MapError(err)
			domain.RenderError(w, status, code, err.Error())
			return
		}
		domain.RenderJSON(w, http.StatusOK, invoices)
		return
	}

	invoices, err := h.svc.ListInvoicesByOrganization(r.Context(), orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, invoices)
}

func (h *BillingHandler) MarkAsPaid(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid invoice ID")
		return
	}

	invoice, err := h.svc.MarkInvoiceAsPaid(r.Context(), id, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, invoice)
}

func (h *BillingHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	var params db.CreateSubscriptionParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}
	params.OrganizationID = orgID

	sub, err := h.svc.CreateSubscription(r.Context(), params)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, sub)
}

func (h *BillingHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	memberIDStr := chi.URLParam(r, "memberID")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid member ID")
		return
	}

	subs, err := h.svc.ListSubscriptionsByMember(r.Context(), memberID, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, subs)
}
