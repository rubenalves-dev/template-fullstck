package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rubenalves-dev/beheer/internal/domain"
)

func withUserClaims(req *http.Request, orgID string) *http.Request {
	claims := &domain.UserClaims{
		OrganizationID: orgID,
		UserID:         uuid.New().String(),
	}
	ctx := context.WithValue(req.Context(), domain.UserClaimsKey, claims)
	return req.WithContext(ctx)
}
