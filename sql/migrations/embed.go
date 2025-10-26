package migrations

import (
	"embed"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
)

//go:embed *.sql
var Embed embed.FS

func MustProvide(pool *pgxpool.Pool) *goose.Provider {
	db := stdlib.OpenDBFromPool(pool)
	provider, err := goose.NewProvider(database.DialectPostgres, db, Embed)
	if err != nil {
		log.Fatalf("migrations provider failed %v", err)
	}
	return provider
}
