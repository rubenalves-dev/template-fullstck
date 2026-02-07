package domain

import "github.com/golang-jwt/jwt/v5"

type ContextKey string

const UserClaimsKey ContextKey = "user_claims"

// UserClaims represents the claims of a JWT token issued to a user.
type UserClaims struct {
	UserID string
	jwt.RegisteredClaims
}
