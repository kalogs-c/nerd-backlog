package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	"github.com/kalogs-c/nerd-backlog/pkg/httpjson"
)

type contextKey string

const accountIDKey contextKey = "account_id"

func AccountIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value := ctx.Value(accountIDKey)
	if value == nil {
		return uuid.Nil, false
	}

	accountID, ok := value.(uuid.UUID)
	return accountID, ok
}

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

			ctx := context.WithValue(r.Context(), accountIDKey, accountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
