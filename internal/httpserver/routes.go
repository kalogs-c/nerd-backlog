package httpserver

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kalogs-c/nerd-backlog/internal/accounts"
	"github.com/kalogs-c/nerd-backlog/internal/games"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

func setupRoutes(
	router chi.Router,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	jwtManager := auth.NewJWTManager([]byte("secret"), time.Minute*5, time.Hour*24)

	router.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(WithAuth(jwtManager, logger))
			setupGames(r, logger, queries)
			setupAccountsProtected(r, logger, queries, jwtManager)
		})

		setupAccounts(r, logger, queries, jwtManager)
	})
}

func setupGames(
	router chi.Router,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	repo := games.NewRepository(queries)
	service := games.NewService(repo)
	adapter := games.NewHTTPAdapter(service, logger)

	router.Get("/games", adapter.ListGames)
	router.Get("/games/{id}", adapter.GetGameByID)
	router.Post("/games", adapter.CreateGame)
	router.Delete("/games/{id}", adapter.DeleteGameByID)
}

func setupAccounts(
	router chi.Router,
	logger *slog.Logger,
	queries *sqlc.Queries,
	jwtManager auth.JWTManager,
) {
	repo := accounts.NewRepository(queries)
	service := accounts.NewService(repo, jwtManager)
	adapter := accounts.NewHTTPAdapter(service, logger)

	router.Post("/login", adapter.Login)
	router.Post("/signup", adapter.Signup)
}

func setupAccountsProtected(
	router chi.Router,
	logger *slog.Logger,
	queries *sqlc.Queries,
	jwtManager auth.JWTManager,
) {
	repo := accounts.NewRepository(queries)
	service := accounts.NewService(repo, jwtManager)
	adapter := accounts.NewHTTPAdapter(service, logger)

	router.Post("/logout", adapter.Logout)
}
