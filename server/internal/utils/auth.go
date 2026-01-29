package domain

import "github.com/golang-jwt/jwt/v5"

type ContextKey string

const (
	UserClaimsKey ContextKey = "user_claims"
)

type UserClaims struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	jwt.RegisteredClaims
}
