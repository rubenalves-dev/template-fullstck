package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rubenalves-dev/beheer/internal/config"
	"github.com/rubenalves-dev/beheer/internal/db"
	"github.com/rubenalves-dev/beheer/internal/auth"
	"github.com/rubenalves-dev/beheer/internal/membership"
	"github.com/rubenalves-dev/beheer/internal/modality"
	"github.com/rubenalves-dev/beheer/internal/billing"
	"github.com/rubenalves-dev/beheer/internal/enrollment"
	httpPort "github.com/rubenalves-dev/beheer/internal/ports/http"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// 2. Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("starting application", "env", cfg.Env, "port", cfg.Port)

	// 3. Database Connection
	pool, err := db.NewPool(context.Background(), cfg.DBConnString)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("connected to database")

	// 4. Repositories & Services
	queries := db.New(pool)
	
	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret)

	memberRepo := membership.NewRepository(queries)
	memberSvc := membership.NewService(memberRepo)

	modalityRepo := modality.NewRepository(queries)
	modalitySvc := modality.NewService(modalityRepo)

	billingRepo := billing.NewRepository(queries)
	billingSvc := billing.NewService(billingRepo)

	enrollmentRepo := enrollment.NewRepository(queries)
	enrollmentSvc := enrollment.NewService(enrollmentRepo, modalitySvc, billingSvc)

	// 5. Setup Router
	r := httpPort.NewRouter(cfg, authSvc, memberSvc, modalitySvc, billingSvc, enrollmentSvc)

	// 6. Start Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// Graceful Shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("server started")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited properly")
}