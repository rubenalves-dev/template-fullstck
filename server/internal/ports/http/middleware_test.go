package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rubenalves-dev/beheer/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret"
	middleware := AuthMiddleware(secret)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(domain.UserClaimsKey).(*domain.UserClaims)
		if ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(claims.UserID))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		userID := uuid.New().String()
		orgID := uuid.New().String()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.UserClaims{
			UserID:         userID,
			OrganizationID: orgID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})
		tokenString, _ := token.SignedString([]byte(secret))

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		middleware(nextHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, userID, rr.Body.String())
	})

	t.Run("missing token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		middleware(nextHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()

		middleware(nextHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}