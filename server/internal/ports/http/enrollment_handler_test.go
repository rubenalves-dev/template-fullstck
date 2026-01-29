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

func TestEnrollmentHandler(t *testing.T) {
	mockSvc := new(mocks.EnrollmentService)
	handler := NewEnrollmentHandler(mockSvc)
	orgID := uuid.New()

	t.Run("Enroll success", func(t *testing.T) {
		reqBody := struct {
			MemberID   uuid.UUID `json:"member_id"`
			ModalityID uuid.UUID `json:"modality_id"`
			SeasonID   uuid.UUID `json:"season_id"`
		}{MemberID: uuid.New(), ModalityID: uuid.New(), SeasonID: uuid.New()}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("EnrollMember", mock.Anything, reqBody.MemberID, reqBody.ModalityID, reqBody.SeasonID, orgID).Return(db.Enrollment{ID: uuid.New()}, nil).Once()

		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.Enroll(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Approve success", func(t *testing.T) {
		enrollmentID := uuid.New()
		mockSvc.On("ApproveEnrollment", mock.Anything, enrollmentID, orgID).Return(db.Enrollment{ID: enrollmentID}, nil).Once()

		req, _ := http.NewRequest("POST", "/"+enrollmentID.String()+"/approve", nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Post("/{id}/approve", handler.Approve)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})
}