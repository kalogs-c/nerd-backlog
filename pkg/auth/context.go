package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	accountIDKey    contextKey = "account_id"
	sessionTokenKey contextKey = "session_token"
)

func WithAccountID(ctx context.Context, accountID uuid.UUID) context.Context {
	return context.WithValue(ctx, accountIDKey, accountID)
}

func AccountIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value := ctx.Value(accountIDKey)
	if value == nil {
		return uuid.Nil, false
	}

	accountID, ok := value.(uuid.UUID)
	return accountID, ok
}

func WithSessionToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, sessionTokenKey, token)
}

func SessionTokenFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(sessionTokenKey)
	if value == nil {
		return "", false
	}

	token, ok := value.(string)
	return token, ok
}
