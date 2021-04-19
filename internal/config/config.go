package config

import (
	"github.com/caarlos0/env/v6"
)

// Config ...
type Config struct {
	ServiceAddress string `env:"SERVICE_ADDRESS" envDefault:"0.0.0.0:8082"`
	DatabaseURI    string `env:"DATABASE_URI,file" envDefault:"./env/database_uri"`
	DatabaseName   string `env:"DATABASE_NAME,file" envDefault:"./env/database_name"`
	GoogleAPIKey   string `env:"GOOGLE_API_KEY,file" envDefault:"./env/google_map_api_key"`
}

// NewConfig ...
func NewConfig() *Config {
	c := new(Config)
	if err := env.Parse(c); err != nil {
		panic(err)
	}
	return c
}
