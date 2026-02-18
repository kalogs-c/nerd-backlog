package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrSessionNotFound = errors.New("session not found")

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
	CreateSession(ctx context.Context, accountID uuid.UUID, token string, expiresAt time.Time) error
	GetSessionAccountID(ctx context.Context, token string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, token string) error
}

type AccountService interface {
	Login(ctx context.Context, email string, password string) (Account, Session, error)
	Register(ctx context.Context, nickname string, email string, password string) (Account, Session, error)
	LogoutSession(ctx context.Context, token string) error
}

type Session struct {
	Token     string
	ExpiresAt time.Time
}
