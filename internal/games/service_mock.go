package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockGameService struct {
	mock.Mock
}

func NewMockGameService() domain.GameService {
	return new(MockGameService)
}

func (m *MockGameService) CreateGame(ctx context.Context, title string) (domain.Game, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(domain.Game), args.Error(1)
}

func (m *MockGameService) GetGameByID(ctx context.Context, id uuid.UUID) (domain.Game, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Game), args.Error(1)
}

func (m *MockGameService) ListGames(ctx context.Context) ([]domain.Game, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Game), args.Error(1)
}

func (m *MockGameService) DeleteGameByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
