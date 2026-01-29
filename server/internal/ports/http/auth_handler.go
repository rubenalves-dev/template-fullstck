package http

import (
	"encoding/json"
	"net/http"

	"github.com/rubenalves-dev/beheer/internal/auth"
	"github.com/rubenalves-dev/beheer/internal/domain"
)

type AuthHandler struct {
	svc auth.Service
}

func NewAuthHandler(svc auth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	FullName         string `json:"full_name"`
	OrganizationName string `json:"organization_name"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusOK, LoginResponse{Token: token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	err := h.svc.Register(r.Context(), auth.RegisterParams{
		Email:            req.Email,
		Password:         req.Password,
		FullName:         req.FullName,
		OrganizationName: req.OrganizationName,
	})
	if err != nil {
		status, code := domain.MapError(err)
		domain.RenderError(w, status, code, err.Error())
		return
	}

	domain.RenderJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}
