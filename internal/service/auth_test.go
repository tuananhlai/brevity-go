package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/repository/mocks"
	"github.com/tuananhlai/brevity-go/internal/service"
)

// hashed value of "password"
var hashedPassword = []byte("$2a$12$lyFRcGsCGdIPv87lZzPn/egx1Nj1xIz6AL628t6auYxHGB9YYYxqW")

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

type AuthServiceTestSuite struct {
	suite.Suite
	authService service.AuthService
	mockRepo    *mocks.AuthRepository
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.mockRepo = mocks.NewAuthRepository(s.T())
	s.authService = service.NewAuthService(s.mockRepo)
}

func (s *AuthServiceTestSuite) TestRegister_Success() {
	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	s.mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(params repository.CreateUserParams) bool {
		return params.Email == email && params.Username == username
	})).Return(&model.AuthUser{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
	}, nil)

	err := s.authService.Register(ctx, email, username, password)

	s.Require().NoError(err)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestRegister_RepositoryError() {
	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	s.mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(params repository.CreateUserParams) bool {
		return params.Email == email && params.Username == username
	})).Return(nil, repository.ErrUserAlreadyExists)

	err := s.authService.Register(ctx, email, username, password)

	s.Require().Error(err)
	s.Require().Equal(repository.ErrUserAlreadyExists, err)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_Success() {
	ctx := context.Background()
	email := "test@example.com"
	password := "password"
	userID := uuid.New()

	s.mockRepo.On("GetUser", ctx, email).Return(&model.AuthUser{
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
	s.Require().NotEmpty(result.RefreshToken)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_UserNotFound() {
	ctx := context.Background()
	email := "nonexistent@example.com"
	password := "password123"

	s.mockRepo.On("GetUser", ctx, email).Return(nil, repository.ErrUserNotFound)

	result, err := s.authService.Login(ctx, email, password)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.Require().Equal(repository.ErrUserNotFound, err)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_InvalidPassword() {
	ctx := context.Background()
	email := "test@example.com"
	password := "wrongpassword"

	s.mockRepo.On("GetUser", ctx, email).Return(&model.AuthUser{
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
