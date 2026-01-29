package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/beheer/internal/domain"
	"github.com/rubenalves-dev/beheer/internal/membership"
)

type MembershipHandler struct {
	svc membership.Service
}

func NewMembershipHandler(svc membership.Service) *MembershipHandler {
	return &MembershipHandler{svc: svc}
}

func (h *MembershipHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	if !ok {
		domain.RenderError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User claims not found")
		return
	}

	orgID, err := uuid.Parse(claims.OrganizationID)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ORG", "Invalid organization ID in token")
		return
	}

	var params membership.CreateMemberParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}
	params.OrganizationID = orgID

	member, err := h.svc.CreateMember(r.Context(), params)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, member)
}

func (h *MembershipHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid member ID")
		return
	}

	member, err := h.svc.GetMember(r.Context(), id, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, member)
}

func (h *MembershipHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	members, err := h.svc.ListMembers(r.Context(), orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, members)
}

func (h *MembershipHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid member ID")
		return
	}

	var params membership.CreateMemberParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}
	params.OrganizationID = orgID

	member, err := h.svc.UpdateMember(r.Context(), id, params)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, member)
}

func (h *MembershipHandler) Archive(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid member ID")
		return
	}

	if err := h.svc.ArchiveMember(r.Context(), id, orgID); err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MembershipHandler) GetFullProfile(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid member ID")
		return
	}

	profile, err := h.svc.GetFullProfile(r.Context(), id, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, profile)
}
