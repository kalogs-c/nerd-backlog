package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const accountIDKey contextKey = "account_id"

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
