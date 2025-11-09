package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockGameRepository struct {
	mock.Mock
}

func NewMockGameRepository() domain.GameRepository {
	return new(MockGameRepository)
}

func (m *MockGameRepository) CreateGame(ctx context.Context, game domain.Game) (domain.Game, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(domain.Game), args.Error(1)
}

func (m *MockGameRepository) GetGameByID(ctx context.Context, id uuid.UUID) (domain.Game, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Game), args.Error(1)
}

func (m *MockGameRepository) ListGames(ctx context.Context) ([]domain.Game, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Game), args.Error(1)
}

func (m *MockGameRepository) DeleteGameByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
