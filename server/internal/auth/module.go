package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/delivery/events"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/delivery/http"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/repositories"
	"github.com/rubenalves-dev/template-fullstack/server/internal/auth/service"
)

type AuthModule struct {
	Service Service
}

func NewModule(pool *pgxpool.Pool, nc *nats.Conn, jwtSecret string) *AuthModule {
	repo := repositories.NewPgxRepository(pool)
	svc := service.NewAuthService(repo, nc, jwtSecret)

	events.RegisterListeners(nc, svc)

	return &AuthModule{Service: svc}
}

func (m *AuthModule) RegisterRoutes(r *chi.Mux) {
	http.RegisterHTTPHandlers(r, m.Service)
}
