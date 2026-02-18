package httpserver

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kalogs-c/nerd-backlog/internal/accounts"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/internal/games"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

func setupRoutes(
	router chi.Router,
	logger *slog.Logger,
	queries *sqlc.Queries,
) {
	sessionManager := auth.NewSessionManager(time.Hour * 24 * 7)
	accountsRepo := accounts.NewRepository(queries)

	router.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(WithAuth(accountsRepo, logger))
			setupGames(r, logger, queries)
			setupAccountsProtected(r, logger, accountsRepo, sessionManager)
		})

		setupAccounts(r, logger, accountsRepo, sessionManager)
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
	repo domain.AccountRepository,
	sessionManager auth.SessionManager,
) {
	service := accounts.NewService(repo, sessionManager)
	adapter := accounts.NewHTTPAdapter(service, logger)

	router.Post("/login", adapter.Login)
	router.Post("/register", adapter.Register)
}

func setupAccountsProtected(
	router chi.Router,
	logger *slog.Logger,
	repo domain.AccountRepository,
	sessionManager auth.SessionManager,
) {
	service := accounts.NewService(repo, sessionManager)
	adapter := accounts.NewHTTPAdapter(service, logger)

	router.Post("/logout", adapter.Logout)
}
