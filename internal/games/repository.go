package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

type repository struct {
	db *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) domain.GameRepository {
	return &repository{q}
}

func (r *repository) CreateGame(ctx context.Context, game domain.Game) (domain.Game, error) {
	insertedGame, err := r.db.CreateGame(ctx, game.Title)
	if err != nil {
		return domain.Game{}, err
	}

	return domain.Game{
		ID:    insertedGame.ID,
		Title: insertedGame.Title,
	}, nil
}

func (r *repository) GetGameByID(ctx context.Context, id uuid.UUID) (domain.Game, error) {
	game, err := r.db.GetGame(ctx, id)
	if err != nil {
		return domain.Game{}, err
	}

	return domain.Game{
		ID:    game.ID,
		Title: game.Title,
	}, nil
}

func (r *repository) ListGames(ctx context.Context) ([]domain.Game, error) {
	games, err := r.db.ListGames(ctx)
	if err != nil {
		return nil, err
	}

	gamesList := make([]domain.Game, len(games))
	for i, game := range games {
		gamesList[i] = domain.Game{ID: game.ID, Title: game.Title}
	}

	return gamesList, nil
}

func (r *repository) DeleteGameByID(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteGameByID(ctx, id)
}
