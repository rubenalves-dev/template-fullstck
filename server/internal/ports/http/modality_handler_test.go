package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/beheer/internal/db"
	"github.com/rubenalves-dev/beheer/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestModalityHandler(t *testing.T) {
	mockSvc := new(mocks.ModalityService)
	handler := NewModalityHandler(mockSvc)
	orgID := uuid.New()

	t.Run("CreateModality success", func(t *testing.T) {
		reqBody := struct {
			Name string `json:"name"`
		}{Name: "Karate"}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("CreateModality", mock.Anything, orgID, reqBody.Name, mock.Anything).Return(db.Modality{ID: uuid.New(), Name: "Karate"}, nil).Once()

		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.CreateModality(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ListModalities success", func(t *testing.T) {
		mockSvc.On("ListModalities", mock.Anything, orgID).Return([]db.Modality{}, nil).Once()
		req, _ := http.NewRequest("GET", "/", nil)
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()
		handler.ListModalities(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("CreateSeason success", func(t *testing.T) {
		reqBody := struct {
			Name    string `json:"name"`
			StartAt string `json:"start_at"`
			EndAt   string `json:"end_at"`
		}{Name: "2026", StartAt: "2026-01-01T00:00:00Z", EndAt: "2026-12-31T23:59:59Z"}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("CreateSeason", mock.Anything, orgID, reqBody.Name, reqBody.StartAt, reqBody.EndAt).Return(db.Season{ID: uuid.New(), Name: "2026"}, nil).Once()

		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.CreateSeason(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ListSeasons success", func(t *testing.T) {
		mockSvc.On("ListSeasons", mock.Anything, orgID).Return([]db.Season{}, nil).Once()
		req, _ := http.NewRequest("GET", "/", nil)
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()
		handler.ListSeasons(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("CreatePrice success", func(t *testing.T) {
		params := db.CreateModalityPriceParams{ModalityID: uuid.New(), SeasonID: uuid.New()}
		body, _ := json.Marshal(params)
		mockSvc.On("CreateModalityPrice", mock.Anything, mock.Anything).Return(db.ModalityPrice{}, nil).Once()
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()
		handler.CreatePrice(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("ListPrices success", func(t *testing.T) {
		seasonID := uuid.New()
		mockSvc.On("ListModalityPrices", mock.Anything, seasonID, orgID).Return([]db.ListModalityPricesBySeasonRow{}, nil).Once()

		req, _ := http.NewRequest("GET", "/"+seasonID.String()+"/prices", nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Get("/{seasonID}/prices", handler.ListPrices)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("CreateWeightCategory success", func(t *testing.T) {
		params := db.CreateWeightCategoryParams{ModalityID: uuid.New(), Name: "Lightweight"}
		body, _ := json.Marshal(params)
		mockSvc.On("CreateWeightCategory", mock.Anything, mock.Anything).Return(db.WeightCategory{}, nil).Once()
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()
		handler.CreateWeightCategory(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("ListWeightCategories success", func(t *testing.T) {
		modalityID := uuid.New()
		mockSvc.On("ListWeightCategories", mock.Anything, modalityID, orgID).Return([]db.WeightCategory{}, nil).Once()

		req, _ := http.NewRequest("GET", "/"+modalityID.String()+"/weight-categories", nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Get("/{modalityID}/weight-categories", handler.ListWeightCategories)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})
}