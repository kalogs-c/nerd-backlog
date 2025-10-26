package config

func newDevHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Host: "localhost",
		Port: "42069",
		DSN:  "postgres://postgres:postgres@localhost:5432/nerd_backlog_dev?sslmode=disable",
	}
}
