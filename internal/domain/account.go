package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

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
}

type AccountService interface {
	Login(ctx context.Context, email, password string) (Account, error)
}

var ErrAccountNotFound = errors.New("account not found")
