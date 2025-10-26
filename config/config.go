package config

type Environment int

const (
	Development Environment = iota
)

type HTTPConfig struct {
	Host string
	Port string
	DSN  string
}

func NewHTTPConfig(environment Environment) *HTTPConfig {
	switch environment {
	case Development:
		return newDevHTTPConfig()
	default:
		return newDevHTTPConfig()
	}
}
