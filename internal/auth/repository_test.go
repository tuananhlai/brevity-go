package auth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/auth"
	"github.com/tuananhlai/brevity-go/internal/testutil"
)

func TestAuthRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

type RepositoryTestSuite struct {
	suite.Suite
	dbTestUtil *testutil.DatabaseTestUtil
	authRepo   auth.Repository
}

func (s *RepositoryTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.authRepo = auth.NewRepository(s.dbTestUtil.DB())
}

func (s *RepositoryTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *RepositoryTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	if err != nil {
		fmt.Println("failed to teardown database: ", err)
	}
}

func (s *RepositoryTestSuite) TestCreateUser_Success() {
	email := "test@test.com"
	passwordHash := []byte("passwordHash")
	username := "test"

	newUser, err := s.authRepo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
	})

	s.Require().NoError(err)
	s.Require().NotNil(newUser)
	s.Require().Equal(email, newUser.Email)
	s.Require().Equal(username, newUser.Username)
}

func (s *RepositoryTestSuite) TestCreateUser_DuplicateEmail() {
	var err error

	email := "test@test.com"
	passwordHash := []byte("passwordHash")
	username := "test"

	_, err = s.authRepo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
	})
	s.Require().NoError(err)

	_, err = s.authRepo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        email,
		PasswordHash: []byte("differentPasswordHash"),
		Username:     "differentUsername",
	})
	s.Require().Error(err)
	s.Require().ErrorIs(err, auth.ErrUserAlreadyExists)
}

func (s *RepositoryTestSuite) TestCreateUser_DuplicateUsername() {
	username := "test"

	_, err := s.authRepo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "test@test.com",
		PasswordHash: []byte("passwordHash"),
		Username:     username,
	})
	s.Require().NoError(err)

	_, err = s.authRepo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "differentEmail@test.com",
		PasswordHash: []byte("differentPasswordHash"),
		Username:     username,
	})
	s.Require().Error(err)
	s.Require().ErrorIs(err, auth.ErrUserAlreadyExists)
}

func (s *RepositoryTestSuite) TestGetUser() {
	email := "test@test.com"
	passwordHash := []byte("passwordHash")
	username := "test"
	ctx := context.Background()

	_, err := s.authRepo.CreateUser(ctx, auth.CreateUserParams{
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
		expectedError    error
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
		{
			name:            "get user by email that does not exist",
			emailOrUsername: "nonexistent@test.com",
			expectedError:   auth.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			user, err := s.authRepo.GetUser(ctx, tc.emailOrUsername)
			if tc.expectedError != nil {
				s.Require().Error(err)
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}

			s.Require().NoError(err)
			s.Require().NotNil(user)
			s.Require().Equal(tc.expectedEmail, user.Email)
			s.Require().Equal(tc.expectedUsername, user.Username)
		})
	}
}
