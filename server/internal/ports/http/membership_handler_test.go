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
	"github.com/rubenalves-dev/beheer/internal/membership"
	"github.com/rubenalves-dev/beheer/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMembershipHandler(t *testing.T) {
	mockSvc := new(mocks.MembershipService)
	handler := NewMembershipHandler(mockSvc)
	orgID := uuid.New()

	t.Run("Create success", func(t *testing.T) {
		params := membership.CreateMemberParams{
			FullName: "John Doe",
		}
		body, _ := json.Marshal(params)

		mockSvc.On("CreateMember", mock.Anything, mock.MatchedBy(func(p membership.CreateMemberParams) bool {
			return p.FullName == "John Doe" && p.OrganizationID == orgID
		})).Return(db.Member{ID: uuid.New(), FullName: "John Doe"}, nil).Once()

		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Get success", func(t *testing.T) {
		memberID := uuid.New()
		mockSvc.On("GetMember", mock.Anything, memberID, orgID).Return(db.Member{ID: memberID}, nil).Once()

		req, _ := http.NewRequest("GET", "/"+memberID.String(), nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		
		r := chi.NewRouter()
		r.Get("/{id}", handler.Get)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("List success", func(t *testing.T) {
		mockSvc.On("ListMembers", mock.Anything, orgID).Return([]db.Member{}, nil).Once()

		req, _ := http.NewRequest("GET", "/", nil)
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Update success", func(t *testing.T) {
		memberID := uuid.New()
		params := membership.CreateMemberParams{FullName: "Jane Doe"}
		body, _ := json.Marshal(params)

		mockSvc.On("UpdateMember", mock.Anything, memberID, mock.Anything).Return(db.Member{ID: memberID, FullName: "Jane Doe"}, nil).Once()

		req, _ := http.NewRequest("PUT", "/"+memberID.String(), bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Put("/{id}", handler.Update)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Archive success", func(t *testing.T) {
		memberID := uuid.New()
		mockSvc.On("ArchiveMember", mock.Anything, memberID, orgID).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/"+memberID.String(), nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Delete("/{id}", handler.Archive)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("GetFullProfile success", func(t *testing.T) {
		memberID := uuid.New()
		mockSvc.On("GetFullProfile", mock.Anything, memberID, orgID).Return(db.GetMemberFullProfileRow{}, nil).Once()

		req, _ := http.NewRequest("GET", "/"+memberID.String()+"/profile", nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Get("/{id}/profile", handler.GetFullProfile)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})
}