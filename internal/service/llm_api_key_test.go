package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
)

func TestLLMAPIKeyService(t *testing.T) {
	suite.Run(t, new(LLMAPIKeyServiceTestSuite))
}

type LLMAPIKeyServiceTestSuite struct {
	suite.Suite
	service        service.LLMAPIKeyService
	mockRepo       *repository.MockLLMAPIKeyRepository
	mockEncryption *service.MockEncryptionService
}

func (s *LLMAPIKeyServiceTestSuite) SetupTest() {
	s.mockRepo = repository.NewMockLLMAPIKeyRepository(s.T())
	s.mockEncryption = service.NewMockEncryptionService(s.T())
	s.service = service.NewLLMAPIKeyService(s.mockRepo, s.mockEncryption)
}

func (s *LLMAPIKeyServiceTestSuite) TestListByUserID_Success() {
	ctx := context.Background()
	userID := uuid.New()
	plainKey := "sk-1234567890abcdef"
	encKey := []byte("encrypted")
	createdAt := time.Now()
	mockModel := &model.LLMAPIKey{
		ID:           uuid.New(),
		Name:         "Test Key",
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    createdAt,
	}

	s.mockRepo.On("ListByUserID", ctx, userID).Return([]*model.LLMAPIKey{mockModel}, nil)
	s.mockEncryption.On("Decrypt", encKey).Return([]byte(plainKey), nil)

	result, err := s.service.ListByUserID(ctx, userID)

	s.Require().NoError(err)
	s.Require().Len(result, 1)
	s.Require().Equal(mockModel.ID, result[0].ID)
	s.Require().Equal("sk-1234567", result[0].ValueFirstTen)
	s.Require().Equal("abcdef", result[0].ValueLastSix)
	s.Require().Equal(userID, result[0].UserID)
	s.Require().Equal(createdAt, result[0].CreatedAt)
	s.mockRepo.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}

func (s *LLMAPIKeyServiceTestSuite) TestListByUserID_RepoError() {
	ctx := context.Background()
	userID := uuid.New()
	s.mockRepo.On("ListByUserID", ctx, userID).Return(nil, assert.AnError)

	result, err := s.service.ListByUserID(ctx, userID)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *LLMAPIKeyServiceTestSuite) TestListByUserID_DecryptError() {
	ctx := context.Background()
	userID := uuid.New()
	encKey := []byte("encrypted")
	mockModel := &model.LLMAPIKey{
		ID:           uuid.New(),
		Name:         "Test Key",
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}

	s.mockRepo.On("ListByUserID", ctx, userID).Return([]*model.LLMAPIKey{mockModel}, nil)
	s.mockEncryption.On("Decrypt", encKey).Return(nil, assert.AnError)

	result, err := s.service.ListByUserID(ctx, userID)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockRepo.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}

func (s *LLMAPIKeyServiceTestSuite) TestCreate_Success() {
	ctx := context.Background()
	userID := uuid.New()
	plainKey := "sk-1234567890abcdef"
	encKey := []byte("encrypted")
	createdAt := time.Now()
	params := service.LLMAPIKeyCreateParams{
		Name:   "Test Key",
		Value:  plainKey,
		UserID: userID.String(),
	}
	mockModel := &model.LLMAPIKey{
		ID:           uuid.New(),
		Name:         params.Name,
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    createdAt,
	}

	s.mockEncryption.On("Encrypt", []byte(plainKey)).Return(encKey)
	s.mockRepo.On("Create", ctx, mock.MatchedBy(func(p repository.LLMAPIKeyCreateParams) bool {
		return p.Name == params.Name && string(p.EncryptedKey) == string(encKey) && p.UserID == userID.String()
	})).Return(mockModel, nil)

	result, err := s.service.Create(ctx, params)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Equal(mockModel.ID, result.ID)
	s.Require().Equal(params.Name, result.Name)
	s.Require().Equal("sk-1234567", result.ValueFirstTen)
	s.Require().Equal("abcdef", result.ValueLastSix)
	s.Require().Equal(userID, result.UserID)
	s.Require().Equal(createdAt, result.CreatedAt)
	s.mockRepo.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}

func (s *LLMAPIKeyServiceTestSuite) TestCreate_RepoError() {
	ctx := context.Background()
	userID := uuid.New()
	plainKey := "sk-1234567890abcdef"
	encKey := []byte("encrypted")
	params := service.LLMAPIKeyCreateParams{
		Name:   "Test Key",
		Value:  plainKey,
		UserID: userID.String(),
	}

	s.mockEncryption.On("Encrypt", []byte(plainKey)).Return(encKey)
	s.mockRepo.On("Create", ctx, mock.Anything).Return(nil, assert.AnError)

	result, err := s.service.Create(ctx, params)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockRepo.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}
