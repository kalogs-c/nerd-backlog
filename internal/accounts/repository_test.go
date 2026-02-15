package accounts

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"testing"
	"time"

	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/internal/storage/postgres"
	"github.com/kalogs-c/nerd-backlog/internal/testutils"
	"github.com/kalogs-c/nerd-backlog/pkg/auth"
	"github.com/kalogs-c/nerd-backlog/sql/migrations"
	sqlc "github.com/kalogs-c/nerd-backlog/sql/sqlc_generated"
	"github.com/stretchr/testify/require"
)

var testQueries *sqlc.Queries

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dsn, terminate, err := testutils.StartPostgresContainer(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	db := postgres.MustConnect(ctx, dsn, nil)
	gooseProvider := migrations.MustProvide(db)
	testQueries = sqlc.New(db)

	_, err = gooseProvider.Up(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	exitCode := m.Run()

	if err := terminate(context.Background()); err != nil {
		log.Println(err)
	}

	os.Exit(exitCode)
}

func TestRepository_CreateAndGetAccount(t *testing.T) {
	repo := NewRepository(testQueries)
	ctx := context.Background()

	nickname := "testuser"
	email := fmt.Sprintf("repo_test%d@example.com", rand.Uint64())
	hashedPassword, _ := auth.HashPassword("password123")

	account, err := repo.CreateAccount(ctx, domain.Account{
		Nickname:       nickname,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	require.NoError(t, err)
	require.NotZero(t, account.ID)

	fetchedAccount, err := repo.GetAccountByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, nickname, fetchedAccount.Nickname)
	require.Equal(t, email, fetchedAccount.Email)
	require.Equal(t, account.ID, fetchedAccount.ID)
}
