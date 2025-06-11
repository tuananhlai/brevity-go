package testutil

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	postgresImageName  = "docker.io/library/postgres:17-alpine"
	postgresDB         = "brevity"
	postgresUser       = "postgres"
	postgresPassword   = "postgres"
	migrationSourceURL = "file://../../db/migrations"
)

// DatabaseTestUtil is a utility struct for testing with a real postgres database.
type DatabaseTestUtil struct {
	db        *sqlx.DB
	container *postgres.PostgresContainer
	migrator  *migrate.Migrate
}

func NewDatabaseTestUtil() (*DatabaseTestUtil, error) {
	ctx := context.Background()
	var err error

	container, err := postgres.Run(ctx, postgresImageName,
		postgres.WithUsername(postgresUser),
		postgres.WithPassword(postgresPassword),
		postgres.WithDatabase(postgresDB),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres connection string: %w", err)
	}

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	migrator, err := migrate.New(
		migrationSourceURL,
		connString,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}
	if err := migrator.Up(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DatabaseTestUtil{
		db:        db,
		container: container,
		migrator:  migrator,
	}, nil
}

// DB returns the postgres database connection.
func (d *DatabaseTestUtil) DB() *sqlx.DB {
	return d.db
}

// Reset resets the database to a clean state by truncating all tables.
func (d *DatabaseTestUtil) Reset() error {
	_, err := d.db.Exec(`
	TRUNCATE TABLE llm_api_keys CASCADE;
	TRUNCATE TABLE articles CASCADE;
	TRUNCATE TABLE users CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	return nil
}

// Teardown closes all resources associated with this database test util instance.
func (d *DatabaseTestUtil) Teardown() error {
	ctx := context.Background()

	err := d.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close postgres connection: %w", err)
	}

	return d.container.Terminate(ctx)
}
