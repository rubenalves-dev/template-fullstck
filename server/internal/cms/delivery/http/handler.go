package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/domain"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/httputil"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/jsonutil"
)

type CMSHandler struct {
	svc domain.Service
}

func RegisterHTTPHandlers(r chi.Router, svc domain.Service) {
	h := &CMSHandler{svc: svc}

	r.Route("/pages", func(r chi.Router) {
		r.Post("/", h.CreateDraft)
		r.Get("/{slug}", h.GetBySlug)
		r.Put("/{id}/metadata", h.UpdateMetadata)
		r.Put("/{id}/layout", h.UpdateLayout)
		r.Post("/{id}/publish", h.Publish)
		r.Post("/{id}/archive", h.Archive)
	})
}

func (h *CMSHandler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	var req domain.CreatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	if err := h.svc.CreateDraft(r.Context(), req.Title); err != nil {
		jsonutil.RenderError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusCreated, map[string]string{"message": "Draft created successfully"})
}

func (h *CMSHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Slug is required")
		return
	}

	page, err := h.svc.GetPageBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, httputil.ErrNotFound) {
			jsonutil.RenderError(w, http.StatusNotFound, "NOT_FOUND", "Page not found")
			return
		}
		jsonutil.RenderError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, page)
}

func (h *CMSHandler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid page ID")
		return
	}

	var req domain.PageUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	if err := h.svc.UpdatePageMetadata(r.Context(), id, req); err != nil {
		jsonutil.RenderError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, map[string]string{"message": "Metadata updated successfully"})
}

func (h *CMSHandler) UpdateLayout(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid page ID")
		return
	}

	var req []domain.RowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse request body")
		return
	}

	if err := h.svc.UpdatePageLayout(r.Context(), id, req); err != nil {
		jsonutil.RenderError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, map[string]string{"message": "Layout updated successfully"})
}

func (h *CMSHandler) Publish(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid page ID")
		return
	}

	if err := h.svc.PublishPage(r.Context(), id); err != nil {
		jsonutil.RenderError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, map[string]string{"message": "Page published successfully"})
}

func (h *CMSHandler) Archive(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		jsonutil.RenderError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid page ID")
		return
	}

	if err := h.svc.ArchivePage(r.Context(), id); err != nil {
		jsonutil.RenderError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	jsonutil.RenderJSON(w, http.StatusOK, map[string]string{"message": "Page archived successfully"})
}
