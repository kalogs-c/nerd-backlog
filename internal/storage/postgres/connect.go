package postgres

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type SlogPgxLogger struct {
	logger *slog.Logger
}

func (l *SlogPgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	if msg != "Query" {
		return
	}

	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	var lvl slog.Level
	switch level {
	case tracelog.LogLevelError:
		lvl = slog.LevelError
	case tracelog.LogLevelWarn:
		lvl = slog.LevelWarn
	case tracelog.LogLevelInfo:
		lvl = slog.LevelInfo
	case tracelog.LogLevelDebug:
		lvl = slog.LevelDebug
	default:
		lvl = slog.LevelInfo
	}

	l.logger.LogAttrs(ctx, lvl, msg, attrs...)
}

func MustConnect(ctx context.Context, dsn string, logger *slog.Logger) *pgxpool.Pool {
	db, err := Connect(ctx, dsn, logger)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	return db
}

func Connect(ctx context.Context, dsn string, logger *slog.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config %s failed: %w", dsn, err)
	}

	if logger != nil {
		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   &SlogPgxLogger{logger: logger},
			LogLevel: tracelog.LogLevelInfo,
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("open db %s failed: %w", dsn, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db error: %w", err)
	}

	return pool, nil
}
