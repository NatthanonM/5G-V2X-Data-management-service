package config

import (
	"github.com/caarlos0/env/v6"
)

// Config ...
type Config struct {
	ServiceAddress string `env:"SERVICE_ADDRESS" envDefault:"0.0.0.0:8082"`
	DatabaseURI    string `env:"DATABASE_URI" envDefault:""`
	DatabaseName   string `env:"DATABASE_NAME" envDefault:""`
	GoogleAPIKey   string `env:"GOOGLE_API_KEY" envDefault:""`
}

// NewConfig ...
func NewConfig() *Config {
	c := new(Config)
	if err := env.Parse(c); err != nil {
		panic(err)
	}
	return c
}
