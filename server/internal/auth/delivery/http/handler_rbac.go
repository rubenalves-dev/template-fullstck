package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/domain"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/httputil"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/jsonutil"
)

type createRoleRequest struct {
	Name string `json:"name"`
}

type assignRoleRequest struct {
	RoleID int `json:"role_id"`
}

type addPermissionRequest struct {
	PermissionID string `json:"permission_id"`
}

func (h *AuthHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req createRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	role, err := h.svc.CreateRole(r.Context(), req.Name)
	if err != nil {
		status, code := httputil.MapError(err)
		jsonutil.RenderError(w, status, code, err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusCreated, role)
}

func (h *AuthHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.GetRoles(r.Context())
	if err != nil {
		status, code := httputil.MapError(err)
		jsonutil.RenderError(w, status, code, err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, roles)
}

func (h *AuthHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_UUID", "Invalid User ID")
		return
	}

	var req assignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	err = h.svc.AssignRole(r.Context(), userID, req.RoleID)
	if err != nil {
		status, code := httputil.MapError(err)
		jsonutil.RenderError(w, status, code, err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, map[string]string{"status": "assigned"})
}

func (h *AuthHandler) AddPermissionToRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "roleID")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_ID", "Invalid Role ID")
		return
	}

	var req addPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	err = h.svc.AddPermissionToRole(r.Context(), roleID, req.PermissionID)
	if err != nil {
		status, code := httputil.MapError(err)
		jsonutil.RenderError(w, status, code, err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, map[string]string{"status": "added"})
}

func (h *AuthHandler) GetMyMenu(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
	if !ok {
		jsonutil.RenderError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context")
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		jsonutil.RenderError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid user ID in token")
		return
	}

	menu, err := h.svc.GetMyMenu(r.Context(), userID)
	if err != nil {
		status, code := httputil.MapError(err)
		jsonutil.RenderError(w, status, code, err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, menu)
}
