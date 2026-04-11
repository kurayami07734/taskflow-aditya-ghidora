package utils

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
	Db   struct {
		Port     int    `env:"DB_PORT" envDefault:"5432"`
		Name     string `env:"DB_NAME"`
		Host     string `env:"DB_HOST"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		SslMode  string `env:"DB_SSL_MODE"`
	}
	JwtSecret string `env:"JWT_SECRET"`
}

func (c Config) GetDbUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Db.User, c.Db.Password, c.Db.Host, c.Db.Port, c.Db.Name, c.Db.SslMode)
}

func LoadConfig() Config {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		slog.Error("Failed to parse config", "error", err)
	}

	return cfg
}
