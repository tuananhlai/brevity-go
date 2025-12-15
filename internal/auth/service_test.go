package auth_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/auth"
)

// hashed value of "password"
var hashedPassword = []byte("$2a$12$lyFRcGsCGdIPv87lZzPn/egx1Nj1xIz6AL628t6auYxHGB9YYYxqW")

func TestAuthService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

type ServiceTestSuite struct {
	suite.Suite
	authService auth.Service
	mockRepo    *auth.MockRepository
}

func (s *ServiceTestSuite) SetupTest() {
	s.mockRepo = auth.NewMockRepository(s.T())
	s.authService = auth.NewService(s.mockRepo, "test-secret")
}

func (s *ServiceTestSuite) TestRegister_Success() {
	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	s.mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(params auth.CreateUserParams) bool {
		return params.Email == email && params.Username == username
	})).Return(&auth.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
	}, nil)

	err := s.authService.Register(ctx, email, username, password)

	s.Require().NoError(err)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestRegister_RepositoryError() {
	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	s.mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(params auth.CreateUserParams) bool {
		return params.Email == email && params.Username == username
	})).Return(nil, auth.ErrUserAlreadyExists)

	err := s.authService.Register(ctx, email, username, password)

	s.Require().Error(err)
	s.Require().ErrorIs(err, auth.ErrUserAlreadyExists)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestLogin_Success() {
	ctx := context.Background()
	email := "test@example.com"
	password := "password"
	userID := uuid.New()

	s.mockRepo.On("GetUser", ctx, email).Return(&auth.User{
		ID:           userID,
		Email:        email,
		Username:     "testuser",
		PasswordHash: hashedPassword,
	}, nil)

	result, err := s.authService.Login(ctx, email, password)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Equal(userID.String(), result.ID)
	s.Require().Equal(email, result.Email)
	s.Require().NotEmpty(result.AccessToken)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestLogin_UserNotFound() {
	ctx := context.Background()
	email := "nonexistent@example.com"
	password := "password123"

	s.mockRepo.On("GetUser", ctx, email).Return(nil, auth.ErrUserNotFound)

	result, err := s.authService.Login(ctx, email, password)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.Require().ErrorIs(err, auth.ErrInvalidCredentials)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestLogin_InvalidPassword() {
	ctx := context.Background()
	email := "test@example.com"
	password := "wrongpassword"

	s.mockRepo.On("GetUser", ctx, email).Return(&auth.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     "testuser",
		PasswordHash: hashedPassword,
	}, nil)

	result, err := s.authService.Login(ctx, email, password)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockRepo.AssertExpectations(s.T())
}
