package accounts

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kalogs-c/nerd-backlog/internal/domain"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
)

type repository struct {
	db *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) domain.AccountRepository {
	return &repository{q}
}

func (r *repository) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	insertedAccount, err := r.db.CreateAccount(ctx, sqlc.CreateAccountParams{
		Nickname:       account.Nickname,
		Email:          account.Email,
		HashedPassword: account.HashedPassword,
	})
	if err != nil {
		return domain.Account{}, err
	}

	return domain.Account{
		ID:             insertedAccount.ID,
		Nickname:       insertedAccount.Nickname,
		Email:          insertedAccount.Email,
		HashedPassword: insertedAccount.HashedPassword,
		TimeStamps: domain.TimeStamps{
			InsertedAt: insertedAccount.InsertedAt.Time,
			UpdatedAt:  insertedAccount.UpdatedAt.Time,
			DeletedAt:  insertedAccount.DeletedAt.Time,
		},
	}, nil
}

func (r *repository) GetAccountByEmail(ctx context.Context, email string) (domain.Account, error) {
	account, err := r.db.GetAccountByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return domain.Account{}, domain.ErrAccountNotFound
	} else if err != nil {
		return domain.Account{}, err
	}

	return domain.Account{
		ID:             account.ID,
		Nickname:       account.Nickname,
		Email:          account.Email,
		HashedPassword: account.HashedPassword,
		TimeStamps: domain.TimeStamps{
			InsertedAt: account.InsertedAt.Time,
			UpdatedAt:  account.UpdatedAt.Time,
			DeletedAt:  account.DeletedAt.Time,
		},
	}, nil
}

func (r *repository) CreateSession(ctx context.Context, accountID uuid.UUID, token string, expiresAt time.Time) error {
	return r.db.CreateSession(ctx, sqlc.CreateSessionParams{
		Token:     token,
		AccountID: accountID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
}

func (r *repository) GetSessionAccountID(ctx context.Context, token string) (uuid.UUID, error) {
	accountID, err := r.db.GetSessionAccountID(ctx, token)
	if err == sql.ErrNoRows {
		return uuid.Nil, domain.ErrSessionNotFound
	} else if err != nil {
		return uuid.Nil, err
	}

	return accountID, nil
}

func (r *repository) DeleteSession(ctx context.Context, token string) error {
	return r.db.DeleteSession(ctx, token)
}
