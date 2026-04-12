package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

var (
	sharedContainer *tcpostgres.PostgresContainer
	sharedDB        *sqlx.DB
	sharedConfig    TestDBConfig
	containerOnce   sync.Once
	containerErr    error
)

func SetupPostgres(t *testing.T) *TestDB {
	t.Helper()

	containerOnce.Do(func() {
		ctx := context.Background()

		pgContainer, err := tcpostgres.Run(ctx,
			"postgres:15-alpine",
			tcpostgres.WithDatabase("taskflow_test"),
			tcpostgres.WithUsername("postgres"),
			tcpostgres.WithPassword("testpass"),
			testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
		)
		if err != nil {
			containerErr = fmt.Errorf("failed to start postgres container: %v", err)
			return
		}

		ip, err := pgContainer.Host(ctx)
		if err != nil {
			containerErr = fmt.Errorf("failed to get container host: %v", err)
			return
		}

		port, err := pgContainer.MappedPort(ctx, "5432")
		if err != nil {
			containerErr = fmt.Errorf("failed to get mapped port: %v", err)
			return
		}

		connectionStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			containerErr = fmt.Errorf("failed to get connection string: %v", err)
			return
		}

		db, err := sqlx.Connect("postgres", connectionStr)
		if err != nil {
			containerErr = fmt.Errorf("failed to connect to database: %v", err)
			return
		}

		if err := runMigrations(db); err != nil {
			containerErr = fmt.Errorf("failed to run migrations: %v", err)
			return
		}

		sharedContainer = pgContainer
		sharedDB = db
		sharedConfig = TestDBConfig{
			Host:     ip,
			Port:     int(port.Num()),
			User:     "postgres",
			Password: "testpass",
			Name:     "taskflow_test",
		}
	})

	if containerErr != nil {
		t.Fatalf("failed to setup postgres: %v", containerErr)
	}

	if err := clearTables(sharedDB); err != nil {
		t.Fatalf("failed to clear tables: %v", err)
	}

	return &TestDB{
		DB:     sharedDB,
		Config: sharedConfig,
		Close:  func() {},
	}
}

func clearTables(db *sqlx.DB) error {
	_, err := db.Exec(`
		TRUNCATE TABLE tasks, projects, users RESTART IDENTITY CASCADE;
	`)
	return err
}

func getMigrationsPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "..", "migrations")
}

func runMigrations(db *sqlx.DB) error {
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
