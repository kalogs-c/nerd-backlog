package config

import "os"

func newDevHTTPConfig() *HTTPConfig {
	host := envOrDefault("HTTP_HOST", "localhost")
	port := envOrDefault("HTTP_PORT", "42069")
	defaultDSN := "postgres://postgres:postgres@localhost:5432/nerd_backlog_dev?sslmode=disable"
	dsn := envOrDefault("DATABASE_DSN", "")
	if dsn == "" {
		dsn = envOrDefault("DATABASE_URL", defaultDSN)
	}

	return &HTTPConfig{
		Host: host,
		Port: port,
		DSN:  dsn,
	}
}

func envOrDefault(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok && value != "" {
		return value
	}

	return fallback
}
