package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/testutil"
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
	dbTestUtil *testutil.DatabaseTestUtil
	authRepo   repository.AuthRepository
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.authRepo = repository.NewAuthRepository(s.dbTestUtil.DB())
}

func (s *AuthRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	if err != nil {
		fmt.Println("failed to teardown database: ", err)
	}
}

func (s *AuthRepositoryTestSuite) TestCreateUser() {
	email := "test@test.com"
	passwordHash := []byte("passwordHash")
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

func (s *AuthRepositoryTestSuite) TestGetUser() {
	email := "test@test.com"
	passwordHash := []byte("passwordHash")
	username := "test"
	ctx := context.Background()

	_, err := s.authRepo.CreateUser(ctx, repository.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
	})
	s.Require().NoError(err)

	testCases := []struct {
		name             string
		emailOrUsername  string
		expectedEmail    string
		expectedUsername string
	}{
		{
			name:             "get user by email",
			emailOrUsername:  email,
			expectedEmail:    email,
			expectedUsername: username,
		},
		{
			name:             "get user by username",
			emailOrUsername:  username,
			expectedEmail:    email,
			expectedUsername: username,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			user, err := s.authRepo.GetUser(ctx, tc.emailOrUsername)
			s.Require().NoError(err)
			s.Require().NotNil(user)
			s.Require().Equal(tc.expectedEmail, user.Email)
			s.Require().Equal(tc.expectedUsername, user.Username)
		})
	}
}
