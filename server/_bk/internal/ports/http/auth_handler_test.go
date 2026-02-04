package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rubenalves-dev/template-fullstack/server/internal/applications/auth"
	"github.com/rubenalves-dev/template-fullstack/server/internal/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		handler := NewAuthHandler(mockSvc)
		
		reqBody := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return("fake-token", nil)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp struct {
			Data LoginResponse `json:"data"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "fake-token", resp.Data.Token)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid request", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		handler := NewAuthHandler(mockSvc)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		handler := NewAuthHandler(mockSvc)

		reqBody := LoginRequest{
			Email:    "test@example.com",
			Password: "wrong",
		}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return("", common.ErrUnauthorized)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		handler := NewAuthHandler(mockSvc)

		reqBody := RegisterRequest{
			Email:            "new@example.com",
			Password:         "password123",
			FullName:         "New User",
			OrganizationName: "New Org",
		}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("Register", mock.Anything, auth.RegisterParams{
			Email:            reqBody.Email,
			Password:         reqBody.Password,
			FullName:         reqBody.FullName,
			OrganizationName: reqBody.OrganizationName,
		}).Return(nil)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Register(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid request", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		handler := NewAuthHandler(mockSvc)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
		rr := httptest.NewRecorder()

		handler.Register(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		handler := NewAuthHandler(mockSvc)

		reqBody := RegisterRequest{
			Email:            "conflict@example.com",
			Password:         "password123",
			FullName:         "Conflict User",
			OrganizationName: "Conflict Org",
		}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("Register", mock.Anything, mock.Anything).Return(common.ErrConflict)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Register(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unexpected service error", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		handler := NewAuthHandler(mockSvc)

		reqBody := RegisterRequest{
			Email:            "error@example.com",
			Password:         "password123",
			FullName:         "Error User",
			OrganizationName: "Error Org",
		}
		body, _ := json.Marshal(reqBody)

		mockSvc.On("Register", mock.Anything, mock.Anything).Return(errors.New("db error"))

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Register(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockSvc.AssertExpectations(t)
	})
}
