package accounts

import (
	"context"

	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
)

type service struct {
	repository domain.AccountRepository
	jwtManager auth.JWTManager
}

func NewService(repository domain.AccountRepository, jwtManager auth.JWTManager) domain.AccountService {
	return &service{repository, jwtManager}
}

func (s *service) Login(ctx context.Context, email string, password string) (domain.Account, domain.TokenPair, error) {
	user, err := s.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	ok, err := auth.ComparePassword(password, user.HashedPassword)
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}
	if !ok {
		return domain.Account{}, domain.TokenPair{}, domain.ErrAccountNotFound
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	refreshToken, exp, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	if err := s.repository.StoreRefreshToken(ctx, user.ID, refreshToken, exp); err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	tokenPair := domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return user, tokenPair, nil
}
