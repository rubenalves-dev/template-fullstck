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

func TestBillingHandler(t *testing.T) {
	mockSvc := new(mocks.BillingService)
	handler := NewBillingHandler(mockSvc)
	orgID := uuid.New()

	t.Run("CreateInvoice success", func(t *testing.T) {
		reqBody := struct {
			Invoice db.CreateInvoiceParams       `json:"invoice"`
			Items   []db.CreateInvoiceItemParams `json:"items"`
		}{
			Invoice: db.CreateInvoiceParams{MemberID: uuid.New()},
			Items:   []db.CreateInvoiceItemParams{{Description: "Item 1"}},
		}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("CreateInvoice", mock.Anything, mock.MatchedBy(func(p db.CreateInvoiceParams) bool {
			return p.OrganizationID == orgID
		}), reqBody.Items).Return(db.Invoice{ID: uuid.New()}, []db.InvoiceItem{}, nil).Once()

		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.CreateInvoice(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("GetInvoice success", func(t *testing.T) {
		invoiceID := uuid.New()
		mockSvc.On("GetInvoice", mock.Anything, invoiceID, orgID).Return(db.Invoice{ID: invoiceID}, []db.InvoiceItem{}, nil).Once()

		req, _ := http.NewRequest("GET", "/"+invoiceID.String(), nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Get("/{id}", handler.GetInvoice)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ListInvoices success", func(t *testing.T) {
		mockSvc.On("ListInvoicesByOrganization", mock.Anything, orgID).Return([]db.ListInvoicesByOrganizationRow{}, nil).Once()

		req, _ := http.NewRequest("GET", "/", nil)
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.ListInvoices(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ListInvoicesByMember success", func(t *testing.T) {
		memberID := uuid.New()
		mockSvc.On("ListInvoicesByMember", mock.Anything, memberID, orgID).Return([]db.Invoice{}, nil).Once()

		req, _ := http.NewRequest("GET", "/?member_id="+memberID.String(), nil)
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()

		handler.ListInvoices(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("MarkAsPaid success", func(t *testing.T) {
		invoiceID := uuid.New()
		mockSvc.On("MarkInvoiceAsPaid", mock.Anything, invoiceID, orgID).Return(db.Invoice{ID: invoiceID}, nil).Once()

		req, _ := http.NewRequest("POST", "/"+invoiceID.String()+"/pay", nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Post("/{id}/pay", handler.MarkAsPaid)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("CreateSubscription success", func(t *testing.T) {
		params := db.CreateSubscriptionParams{MemberID: uuid.New(), ModalityPriceID: uuid.New()}
		body, _ := json.Marshal(params)
		mockSvc.On("CreateSubscription", mock.Anything, mock.Anything).Return(db.Subscription{}, nil).Once()
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req = withUserClaims(req, orgID.String())
		rr := httptest.NewRecorder()
		handler.CreateSubscription(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("ListSubscriptions success", func(t *testing.T) {
		memberID := uuid.New()
		mockSvc.On("ListSubscriptionsByMember", mock.Anything, memberID, orgID).Return([]db.ListSubscriptionsByMemberRow{}, nil).Once()

		req, _ := http.NewRequest("GET", "/member/"+memberID.String(), nil)
		req = withUserClaims(req, orgID.String())
		
		rr := httptest.NewRecorder()
		r := chi.NewRouter()
		r.Get("/member/{memberID}", handler.ListSubscriptions)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockSvc.AssertExpectations(t)
	})
}
