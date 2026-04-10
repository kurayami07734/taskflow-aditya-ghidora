package main

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to read .env file: %v", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return cfg
}
