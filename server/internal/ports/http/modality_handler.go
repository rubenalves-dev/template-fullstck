package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/beheer/internal/db"
	"github.com/rubenalves-dev/beheer/internal/domain"
	"github.com/rubenalves-dev/beheer/internal/modality"
)

type ModalityHandler struct {
	svc modality.Service
}

func NewModalityHandler(svc modality.Service) *ModalityHandler {
	return &ModalityHandler{svc: svc}
}

func (h *ModalityHandler) CreateModality(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	var req struct {
		Name  string  `json:"name"`
		Notes *string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	m, err := h.svc.CreateModality(r.Context(), orgID, req.Name, req.Notes)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, m)
}

func (h *ModalityHandler) ListModalities(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	modalities, err := h.svc.ListModalities(r.Context(), orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, modalities)
}

func (h *ModalityHandler) CreateSeason(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	var req struct {
		Name    string `json:"name"`
		StartAt string `json:"start_at"`
		EndAt   string `json:"end_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	s, err := h.svc.CreateSeason(r.Context(), orgID, req.Name, req.StartAt, req.EndAt)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, s)
}

func (h *ModalityHandler) ListSeasons(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	seasons, err := h.svc.ListSeasons(r.Context(), orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, seasons)
}

func (h *ModalityHandler) CreatePrice(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	var params db.CreateModalityPriceParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}
	params.OrganizationID = orgID

	p, err := h.svc.CreateModalityPrice(r.Context(), params)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, p)
}

func (h *ModalityHandler) ListPrices(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	seasonIDStr := chi.URLParam(r, "seasonID")
	seasonID, err := uuid.Parse(seasonIDStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid season ID")
		return
	}

	prices, err := h.svc.ListModalityPrices(r.Context(), seasonID, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, prices)
}

func (h *ModalityHandler) CreateWeightCategory(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	var params db.CreateWeightCategoryParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}
	params.OrganizationID = orgID

	wc, err := h.svc.CreateWeightCategory(r.Context(), params)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, wc)
}

func (h *ModalityHandler) ListWeightCategories(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	modalityIDStr := chi.URLParam(r, "modalityID")
	modalityID, err := uuid.Parse(modalityIDStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid modality ID")
		return
	}

	categories, err := h.svc.ListWeightCategories(r.Context(), modalityID, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, categories)
}
