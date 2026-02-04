package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/rubenalves-dev/template-fullstack/server/internal/applications/auth"
	"github.com/rubenalves-dev/template-fullstack/server/internal/platform"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/jsonutil"
)

func NewRouter(cfg *platform.Config, authSvc auth.Service) http.Handler {
	r := chi.NewRouter()

	// Standard Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // In production, we might want a custom slog middleware
	r.Use(middleware.Recoverer)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust as needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(c.Handler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		jsonutil.RenderJSON(w, http.StatusOK, map[string]string{"status": "OK"})
	})

	// Auth
	authHandler := NewAuthHandler(authSvc)
	r.Post("/login", authHandler.Login)
	r.Post("/register", authHandler.Register)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(cfg.JWTSecret))

		/*
			// Membership
			memberHandler := NewMembershipHandler(memberSvc)
			r.Route("/api/v1/members", func(r chi.Router) {
				r.Post("/", memberHandler.Create)
				r.Get("/", memberHandler.List)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", memberHandler.Get)
					r.Put("/", memberHandler.Update)
					r.Delete("/", memberHandler.Archive)
					r.Get("/profile", memberHandler.GetFullProfile)
				})
			})

			// Modalities & Seasons
			modalityHandler := NewModalityHandler(modalitySvc)
			r.Route("/api/v1/modalities", func(r chi.Router) {
				r.Post("/", modalityHandler.CreateModality)
				r.Get("/", modalityHandler.ListModalities)
				r.Get("/{modalityID}/weight-categories", modalityHandler.ListWeightCategories)
				r.Post("/weight-categories", modalityHandler.CreateWeightCategory)
			})
			r.Route("/api/v1/seasons", func(r chi.Router) {
				r.Post("/", modalityHandler.CreateSeason)
				r.Get("/", modalityHandler.ListSeasons)
				r.Get("/{seasonID}/prices", modalityHandler.ListPrices)
				r.Post("/prices", modalityHandler.CreatePrice)
			})

			// Billing
			billingHandler := NewBillingHandler(billingSvc)
			r.Route("/api/v1/billing", func(r chi.Router) {
				r.Route("/invoices", func(r chi.Router) {
					r.Post("/", billingHandler.CreateInvoice)
					r.Get("/", billingHandler.ListInvoices)
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", billingHandler.GetInvoice)
						r.Post("/pay", billingHandler.MarkAsPaid)
					})
				})
				r.Route("/subscriptions", func(r chi.Router) {
					r.Post("/", billingHandler.CreateSubscription)
					r.Get("/member/{memberID}", billingHandler.ListSubscriptions)
				})
			})

			// Enrollment
			enrollmentHandler := NewEnrollmentHandler(enrollmentSvc)
			r.Route("/api/v1/enrollments", func(r chi.Router) {
				r.Post("/", enrollmentHandler.Enroll)
				r.Get("/member/{memberID}", enrollmentHandler.ListByMember)
				r.Post("/{id}/approve", enrollmentHandler.Approve)
			})
		*/
	})

	return r
}
