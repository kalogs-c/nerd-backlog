package accounts

import (
	"context"
	"database/sql"

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
