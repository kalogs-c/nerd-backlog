package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
)

type service struct {
	repository     domain.AccountRepository
	sessionManager auth.SessionManager
}

func NewService(repository domain.AccountRepository, sessionManager auth.SessionManager) domain.AccountService {
	return &service{repository, sessionManager}
}

func (s *service) Login(ctx context.Context, email string, password string) (domain.Account, domain.Session, error) {
	user, err := s.authenticate(ctx, email, password)
	if err != nil {
		return domain.Account{}, domain.Session{}, err
	}

	session, err := s.issueSession(ctx, user.ID)
	if err != nil {
		return domain.Account{}, domain.Session{}, err
	}

	return user, session, nil
}

func (s *service) Register(ctx context.Context, nickname string, email string, password string) (domain.Account, domain.Session, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return domain.Account{}, domain.Session{}, err
	}

	account, err := s.repository.CreateAccount(ctx, domain.Account{
		Nickname:       nickname,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return domain.Account{}, domain.Session{}, err
	}

	session, err := s.issueSession(ctx, account.ID)
	if err != nil {
		return domain.Account{}, domain.Session{}, err
	}

	return account, session, nil
}

func (s *service) LogoutSession(ctx context.Context, token string) error {
	return s.repository.DeleteSession(ctx, token)
}

func (s *service) authenticate(ctx context.Context, email string, password string) (domain.Account, error) {
	user, err := s.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, err
	}

	ok, err := auth.ComparePassword(password, user.HashedPassword)
	if err != nil {
		return domain.Account{}, err
	}
	if !ok {
		return domain.Account{}, domain.ErrAccountNotFound
	}

	return user, nil
}

func (s *service) issueSession(ctx context.Context, accountID uuid.UUID) (domain.Session, error) {
	token, expiresAt, err := s.sessionManager.GenerateSessionToken()
	if err != nil {
		return domain.Session{}, err
	}

	if err := s.repository.CreateSession(ctx, accountID, token, expiresAt); err != nil {
		return domain.Session{}, err
	}

	return domain.Session{Token: token, ExpiresAt: expiresAt}, nil
}
