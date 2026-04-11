package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port       int    `env:"PORT" envDefault:"8080"`
	DbPort     int    `env:"DB_PORT" envDefault:"5432"`
	DbName     string `env:"DB_NAME"`
	DbHost     string `env:"DB_HOST"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbSslMode  string `env:"DB_SSL_MODE"`
}

func (c Config) getDbUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbSslMode)
}

func LoadConfig() Config {
	// NOTE: Error is ignored because the .env file will notexist inside the container
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return cfg
}
