package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/pkg/validator"
)

type GameResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"topic"`
}

type CreateGamePayload struct {
	Title string `json:"topic"`
}

func (crp *CreateGamePayload) Valid(ctx context.Context) validator.Problems {
	problems := make(validator.Problems)
	return problems
}
