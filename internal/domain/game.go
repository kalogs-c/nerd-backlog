package domain

import (
	"context"

	"github.com/google/uuid"
)

type Game struct {
	ID    uuid.UUID
	Title string
}

type GameService interface {
	CreateGame(ctx context.Context, title string) (Game, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (Game, error)
	ListGames(ctx context.Context) ([]Game, error)
}

type GameRepository interface {
	CreateGame(ctx context.Context, game Game) (Game, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (Game, error)
	ListGames(ctx context.Context) ([]Game, error)
}
