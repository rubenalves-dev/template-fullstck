package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	iam_auth "github.com/rubenalves-dev/template-fullstack/server/internal/domain/iam/auth"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/jsonutil"
)

// AuthMiddleware validates the JWT token and injects UserClaims into the context.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				jsonutil.RenderError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				jsonutil.RenderError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid authorization header")
				return
			}

			tokenString := parts[1]
			claims := &iam_auth.UserClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				jsonutil.RenderError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), iam_auth.UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}