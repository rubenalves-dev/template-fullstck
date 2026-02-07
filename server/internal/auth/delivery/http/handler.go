package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/domain"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/httputil"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/jsonutil"
)

type AuthHandler struct {
	svc domain.Service
}

func RegisterHTTPHandlers(r *chi.Mux, svc domain.Service) {
	h := &AuthHandler{svc: svc}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})
}

func RegisterProtectedHTTPHandlers(r chi.Router, svc domain.Service) {
	h := &AuthHandler{svc: svc}

	r.Route("/backoffice", func(r chi.Router) {
		r.Get("/me/menu", h.GetMyMenu)
		r.Get("/roles", h.GetRoles)
		r.Post("/roles", h.CreateRole)
		r.Post("/roles/{roleID}/permissions", h.AddPermissionToRole)
		r.Post("/users/{userID}/roles", h.AssignRoleToUser)
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		status, code := httputil.MapError(err)
		jsonutil.RenderError(w, status, code, err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	err := h.svc.Register(r.Context(), domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		FullName:     req.FullName,
	})

	if err != nil {
		status, code := httputil.MapError(err)
		jsonutil.RenderError(w, status, code, err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}
