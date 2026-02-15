package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/kalogs-c/nerd-backlog/config"
	"github.com/kalogs-c/nerd-backlog/internal/httpserver"
	"github.com/kalogs-c/nerd-backlog/internal/storage/postgres"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := config.NewHTTPConfig(config.Development)

	db := postgres.MustConnect(ctx, config.DSN, logger)
	queries := sqlc.New(db)

	server := httpserver.NewHTTPServer(
		logger,
		queries,
		config,
		middleware.RequestID,
		middleware.Recoverer,
		middleware.StripSlashes,
		httpserver.WithLogging(logger),
	)

	go server.MustServe()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("Server gracefully stopped")
}
