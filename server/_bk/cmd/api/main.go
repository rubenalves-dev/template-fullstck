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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rubenalves-dev/template-fullstack/server/internal/applications/auth"
	"github.com/rubenalves-dev/template-fullstack/server/internal/infrastructure/persistence/postgres/repositories"
	"github.com/rubenalves-dev/template-fullstack/server/internal/platform"
	httpPort "github.com/rubenalves-dev/template-fullstack/server/internal/ports/http"
)

func main() {
	// 1. Load Configuration
	cfg, err := platform.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// 2. Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("starting application", "env", cfg.Env, "port", cfg.Port)

	// 3. Database Connection
	pool, err := pgxpool.New(context.Background(), cfg.DBConnString)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("connected to database")

	// 4. Repositories & Services
	authRepo := repositories.NewUserRepository(pool)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret)

	// 5. Setup Router
	r := httpPort.NewRouter(cfg, authSvc)

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
