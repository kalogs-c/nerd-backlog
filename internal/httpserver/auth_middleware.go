package httpserver

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	"github.com/kalogs-c/nerd-backlog/pkg/httpjson"
)

type SessionStore interface {
	GetSessionAccountID(ctx context.Context, token string) (uuid.UUID, error)
}

func WithAuth(sessionStore SessionStore, logger *slog.Logger) Middleware {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sessionStore == nil {
				httpjson.NotifyHTTPError(w, r, logger, http.StatusUnauthorized, "missing session", auth.ErrInvalidToken)
				return
			}

			cookie, err := r.Cookie(auth.SessionCookieName)
			if err != nil || cookie.Value == "" {
				httpjson.NotifyHTTPError(w, r, logger, http.StatusUnauthorized, "missing session", auth.ErrInvalidToken)
				return
			}

			accountID, err := sessionStore.GetSessionAccountID(r.Context(), cookie.Value)
			if err != nil {
				httpjson.NotifyHTTPError(w, r, logger, http.StatusUnauthorized, "invalid session", err)
				return
			}

			ctx := auth.WithAccountID(r.Context(), accountID)
			ctx = auth.WithSessionToken(ctx, cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
