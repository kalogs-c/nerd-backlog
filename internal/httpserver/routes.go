package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kalogs-c/nerd-backlog/internal/accounts"
	"github.com/kalogs-c/nerd-backlog/internal/games"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

func setupRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	setupGames(mux, logger, queries)
	setupAccounts(mux, logger, queries)
}

func setupGames(
	mux *http.ServeMux,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	repo := games.NewRepository(queries)
	service := games.NewService(repo)
	adapter := games.NewHTTPAdapter(service, logger)

	mux.HandleFunc("GET /api/games", adapter.ListGames)
	mux.HandleFunc("GET /api/games/{id}", adapter.GetGameByID)
	mux.HandleFunc("POST /api/games", adapter.CreateGame)
	mux.HandleFunc("DELETE /api/games/{id}", adapter.DeleteGameByID)
}

func setupAccounts(
	mux *http.ServeMux,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	repo := accounts.NewRepository(queries)
	service := accounts.NewService(repo, auth.NewJWTManager([]byte("secret"), time.Minute*5, time.Hour*24))
	adapter := accounts.NewHTTPAdapter(service, logger)

	mux.HandleFunc("GET /api/login", adapter.Login)
}
