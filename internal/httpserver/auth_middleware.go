package httpserver

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	"github.com/kalogs-c/nerd-backlog/pkg/httpjson"
)

func WithAuth(jwtManager auth.JWTManager, logger *slog.Logger) Middleware {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, "Bearer ") {
				httpjson.NotifyHTTPError(w, r, logger, http.StatusUnauthorized, "missing bearer token", auth.ErrInvalidToken)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
			if token == "" {
				httpjson.NotifyHTTPError(w, r, logger, http.StatusUnauthorized, "missing bearer token", auth.ErrInvalidToken)
				return
			}

			accountID, err := jwtManager.VerifyAccessToken(token)
			if err != nil {
				httpjson.NotifyHTTPError(w, r, logger, http.StatusUnauthorized, "invalid token", err)
				return
			}

			ctx := auth.WithAccountID(r.Context(), accountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
