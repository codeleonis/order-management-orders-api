package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/your-service/config"
	productrepo "github.com/your-org/your-service/internal/infrastructure/postgres/product"
	createproduct "github.com/your-org/your-service/internal/product/create"
	deleteproduct "github.com/your-org/your-service/internal/product/delete"
	findbyid "github.com/your-org/your-service/internal/product/find_by_id"
	listproducts "github.com/your-org/your-service/internal/product/list"
	updateproduct "github.com/your-org/your-service/internal/product/update"
	"github.com/your-org/your-service/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Infrastructure
	repo := productrepo.NewRepository(pool)

	// Use cases
	createUC := createproduct.New(repo)
	deleteUC := deleteproduct.New(repo)
	findByIDUC := findbyid.New(repo)
	listUC := listproducts.New(repo)
	updateUC := updateproduct.New(repo)

	// Handlers
	handlers := server.Handlers{
		Create:   createproduct.NewHandler(createUC),
		Delete:   deleteproduct.NewHandler(deleteUC),
		FindByID: findbyid.NewHandler(findByIDUC),
		List:     listproducts.NewHandler(listUC),
		Update:   updateproduct.NewHandler(updateUC),
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      server.NewRouter(handlers),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
	slog.Info("server stopped")
}
