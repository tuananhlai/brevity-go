package repository_test

import (
	"context"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/tuananhlai/brevity-go/internal/repository"
)

const (
	postgresImageName = "docker.io/library/postgres:17-alpine"
	postgresDB        = "brevity"
	postgresUser      = "postgres"
	postgresPassword  = "postgres"
)

func TestAuthRepository(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

type AuthRepositoryTestSuite struct {
	suite.Suite
	db        *sqlx.DB
	authRepo  repository.AuthRepository
	container *postgres.PostgresContainer
	migrator  *migrate.Migrate
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	var err error

	s.container, err = postgres.Run(ctx, postgresImageName,
		postgres.WithUsername(postgresUser),
		postgres.WithPassword(postgresPassword),
		postgres.WithDatabase(postgresDB),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	connString, err := s.container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get postgres connection string: %v", err)
	}

	s.db, err = sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	s.migrator, err = migrate.New("file://../../db/migrations", connString)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	if err := s.migrator.Up(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	s.authRepo = repository.NewAuthRepository(s.db)
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	s.db.Close()
	s.container.Terminate(ctx)
}

func (s *AuthRepositoryTestSuite) TestCreateUser() {
	email := "test@test.com"
	passwordHash := "passwordHash"
	username := "test"

	newUser, err := s.authRepo.CreateUser(context.Background(), repository.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
	})

	s.Require().NoError(err)
	s.Require().NotNil(newUser)
	s.Require().Equal(email, newUser.Email)
	s.Require().Equal(username, newUser.Username)
}
