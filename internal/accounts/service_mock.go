package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func NewMockAccountService() domain.AccountService {
	return new(MockAccountService)
}

func (m *MockAccountService) Login(ctx context.Context, email string, password string) (domain.Account, domain.TokenPair, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(domain.Account), args.Get(1).(domain.TokenPair), args.Error(2)
}

func (m *MockAccountService) Register(ctx context.Context, nickname string, email string, password string) (domain.Account, domain.TokenPair, error) {
	args := m.Called(ctx, nickname, email, password)
	return args.Get(0).(domain.Account), args.Get(1).(domain.TokenPair), args.Error(2)
}

func (m *MockAccountService) Logout(ctx context.Context, accountID uuid.UUID) error {
	args := m.Called(ctx, accountID)
	return args.Error(0)
}
