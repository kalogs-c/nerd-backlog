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

func (m *MockAccountRepository) CreateSession(ctx context.Context, accountID uuid.UUID, token string, expiresAt time.Time) error {
	args := m.Called(ctx, accountID, token, expiresAt)
	return args.Error(0)
}

func (m *MockAccountRepository) GetSessionAccountID(ctx context.Context, token string) (uuid.UUID, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAccountRepository) DeleteSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
