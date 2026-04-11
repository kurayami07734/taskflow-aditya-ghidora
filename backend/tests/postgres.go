package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	DB     *sqlx.DB
	Close  func()
	Config TestDBConfig
}

type TestDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func SetupPostgres(t *testing.T) *TestDB {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase("taskflow_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	ip, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	connectionStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sqlx.Connect("postgres", connectionStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := runMigrations(db, t); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	cfg := TestDBConfig{
		Host:     ip,
		Port:     int(port.Num()),
		User:     "postgres",
		Password: "testpass",
		Name:     "taskflow_test",
	}

	return &TestDB{
		DB:     db,
		Config: cfg,
		Close: func() {
			db.Close()
			pgContainer.Terminate(context.Background())
		},
	}
}

func getMigrationsPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "..", "migrations")
}

func runMigrations(db *sqlx.DB, t *testing.T) error {
	migrationsPath := getMigrationsPath()
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, f := range files {
		name := f.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(migrationsPath, name))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", name, err)
		}
	}
	return nil
}

func (cfg TestDBConfig) GetDbUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
