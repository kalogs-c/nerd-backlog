package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrAccountNotFound = errors.New("account not found")

type Account struct {
	ID             uuid.UUID
	Nickname       string
	Email          string
	HashedPassword string
	TimeStamps
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, user Account) (Account, error)
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiresAt time.Time) error
}

type AccountService interface {
	Login(ctx context.Context, email string, password string) (Account, TokenPair, error)
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
