package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
)

type service struct {
	repository domain.GameRepository
}

func NewService(gameRepository domain.GameRepository) domain.GameService {
	return &service{gameRepository}
}

func (s *service) CreateGame(ctx context.Context, title string) (domain.Game, error) {
	game := domain.Game{Title: title}
	return s.repository.CreateGame(ctx, game)
}

func (s *service) GetGameByID(ctx context.Context, id uuid.UUID) (domain.Game, error) {
	return s.repository.GetGameByID(ctx, id)
}

func (s *service) ListGames(ctx context.Context) ([]domain.Game, error) {
	return s.repository.ListGames(ctx)
}

func (s *service) DeleteGameByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.repository.GetGameByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repository.DeleteGameByID(ctx, id)
}
