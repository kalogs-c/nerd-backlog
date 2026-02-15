package accounts

import (
	"context"

	"github.com/google/uuid"
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

func (s *service) Register(ctx context.Context, nickname string, email string, password string) (domain.Account, domain.TokenPair, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	account, err := s.repository.CreateAccount(ctx, domain.Account{
		Nickname:       nickname,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(account.ID)
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	refreshToken, exp, err := s.jwtManager.GenerateRefreshToken(account.ID)
	if err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	if err := s.repository.StoreRefreshToken(ctx, account.ID, refreshToken, exp); err != nil {
		return domain.Account{}, domain.TokenPair{}, err
	}

	tokenPair := domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return account, tokenPair, nil
}

func (s *service) Logout(ctx context.Context, accountID uuid.UUID) error {
	return s.repository.DeleteRefreshToken(ctx, accountID)
}
