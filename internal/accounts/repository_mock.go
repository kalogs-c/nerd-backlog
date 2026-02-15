package accounts

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func NewMockAccountRepository() domain.AccountRepository {
	return new(MockAccountRepository)
}

func (m *MockAccountRepository) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByEmail(ctx context.Context, email string) (domain.Account, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *MockAccountRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, refreshToken, expiresAt)
	return args.Error(0)
}

func (m *MockAccountRepository) DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
