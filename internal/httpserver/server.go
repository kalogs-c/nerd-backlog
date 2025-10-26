package httpserver

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/kalogs-c/nerd-backlog/config"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

type HTTPServer struct {
	logger  *slog.Logger
	queries *sqlc.Queries
	config  *config.HTTPConfig
	server  http.Server
}

func NewHTTPServer(
	logger *slog.Logger,
	queries *sqlc.Queries,
	config *config.HTTPConfig,
	middlewares ...Middleware,
) *HTTPServer {
	mux := http.NewServeMux()

	setupRoutes(mux, logger, queries)

	var handler http.Handler = mux
	for _, m := range middlewares {
		handler = m(handler)
	}

	return &HTTPServer{
		logger:  logger,
		queries: queries,
		config:  config,
		server: http.Server{
			Addr:    net.JoinHostPort(config.Host, config.Port),
			Handler: handler,
		},
	}
}

func (s *HTTPServer) MustServe() {
	s.logger.Info("Server up and running!", "host", s.config.Host, "port", s.config.Port)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("error listening and serving", "err", err)
		os.Exit(1)
	}
}
