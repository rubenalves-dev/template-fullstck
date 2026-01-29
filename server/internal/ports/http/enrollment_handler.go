package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/beheer/internal/domain"
	"github.com/rubenalves-dev/beheer/internal/enrollment"
)

type EnrollmentHandler struct {
	svc enrollment.Service
}

func NewEnrollmentHandler(svc enrollment.Service) *EnrollmentHandler {
	return &EnrollmentHandler{svc: svc}
}

func (h *EnrollmentHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	var req struct {
		MemberID   uuid.UUID `json:"member_id"`
		ModalityID uuid.UUID `json:"modality_id"`
		SeasonID   uuid.UUID `json:"season_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	e, err := h.svc.EnrollMember(r.Context(), req.MemberID, req.ModalityID, req.SeasonID, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, e)
}

func (h *EnrollmentHandler) ListByMember(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	memberIDStr := chi.URLParam(r, "memberID")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid member ID")
		return
	}

	enrollments, err := h.svc.ListEnrollmentsByMember(r.Context(), memberID, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, enrollments)
}

func (h *EnrollmentHandler) Approve(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	orgID, _ := uuid.Parse(claims.OrganizationID)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid enrollment ID")
		return
	}

	e, err := h.svc.ApproveEnrollment(r.Context(), id, orgID)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, e)
}
