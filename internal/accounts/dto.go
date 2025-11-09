package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/pkg/validator"
)

type AccountResponse struct {
	ID       uuid.UUID `json:"id"`
	Nickname string    `json:"nickname"`
	Email    string    `json:"email"`
}

func MountAccountResponse(account domain.Account) AccountResponse {
	return AccountResponse{
		ID:       account.ID,
		Nickname: account.Nickname,
		Email:    account.Email,
	}
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func MountTokenPairResponse(tokenPair domain.TokenPair) TokenPairResponse {
	return TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
}

type LoginResponse struct {
	Account   AccountResponse   `json:"account"`
	TokenPair TokenPairResponse `json:"token_pair"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lp *LoginPayload) Valid(ctx context.Context) validator.Problems {
	problems := make(validator.Problems)
	return problems
}
