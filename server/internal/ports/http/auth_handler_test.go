package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rubenalves-dev/beheer/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	mockSvc := new(mocks.AuthService)
	handler := NewAuthHandler(mockSvc)

	t.Run("success", func(t *testing.T) {
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
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "fake-token", resp.Data.Token)
	})

	t.Run("invalid request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}