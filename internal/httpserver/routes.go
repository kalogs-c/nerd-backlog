package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/kalogs-c/nerd-backlog/internal/games"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

func setupRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	setupGames(mux, logger, queries)
}

func setupGames(
	mux *http.ServeMux,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	repo := games.NewRepository(queries)
	service := games.NewService(repo)
	adapter := games.NewHTTPAdapter(service, logger)

	mux.HandleFunc("GET /games", adapter.ListGames)
	mux.HandleFunc("GET /games/{id}", adapter.GetGameByID)
	mux.HandleFunc("POST /games", adapter.CreateGame)
	mux.HandleFunc("DELETE /games/{id}", adapter.DeleteGameByID)
}
