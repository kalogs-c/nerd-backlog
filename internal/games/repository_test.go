package games

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kalogs-c/nerd-backlog/internal/domain"
	"github.com/kalogs-c/nerd-backlog/internal/storage/postgres"
	"github.com/kalogs-c/nerd-backlog/internal/testutils"
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

func TestRepository_CreateAndGetGame(t *testing.T) {
	repo := NewRepository(testQueries)
	ctx := context.Background()

	game, err := repo.CreateGame(ctx, domain.Game{Title: "Backlog, the game"})
	require.NoError(t, err)
	require.NotZero(t, game.ID)

	got, err := repo.GetGameByID(ctx, game.ID)
	require.NoError(t, err)
	require.Equal(t, "Backlog, the game", got.Title)
	require.Equal(t, game.ID, got.ID)
}

func TestRepository_DeleteGame(t *testing.T) {
	repo := NewRepository(testQueries)
	ctx := context.Background()

	game, err := repo.CreateGame(ctx, domain.Game{Title: "Backlog, the game"})
	require.NoError(t, err)
	require.NotZero(t, game.ID)

	err = repo.DeleteGameByID(ctx, game.ID)
	require.NoError(t, err)
}
