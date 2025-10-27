package domain

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Nickname       string
	Email          string
	HashedPassword string
	TimeStamps
}

type AccountRepository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type AccountService interface {
	Login(ctx context.Context, email, password string) (User, error)
	Register(ctx context.Context, user User) (User, error)
}
