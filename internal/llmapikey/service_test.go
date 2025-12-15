package llmapikey_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/llmapikey"
)

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

type ServiceTestSuite struct {
	suite.Suite
	service        llmapikey.Service
	mockRepo       *llmapikey.MockRepository
	mockEncryption *llmapikey.MockCrypter
}

func (s *ServiceTestSuite) SetupTest() {
	s.mockRepo = llmapikey.NewMockRepository(s.T())
	s.mockEncryption = llmapikey.NewMockCrypter(s.T())
	s.service = llmapikey.NewService(s.mockRepo, s.mockEncryption)
}

func (s *ServiceTestSuite) TestListByUserID_Success() {
	ctx := context.Background()
	userID := uuid.New()
	plainKey := "sk-1234567890abcdef"
	encKey := []byte("encrypted")
	createdAt := time.Now()
	mockModel := &llmapikey.StoredAPIKey{
		ID:           uuid.New(),
		Name:         "Test Key",
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    createdAt,
	}

	s.mockRepo.On("ListByUserID", ctx, userID.String()).Return([]*llmapikey.StoredAPIKey{
		mockModel,
	}, nil)
	s.mockEncryption.On("Decrypt", encKey).Return([]byte(plainKey), nil)

	result, err := s.service.ListByUserID(ctx, userID.String())

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

func (s *ServiceTestSuite) TestListByUserID_RepoError() {
	ctx := context.Background()
	userID := uuid.New()
	s.mockRepo.On("ListByUserID", ctx, userID.String()).Return(nil, assert.AnError)

	result, err := s.service.ListByUserID(ctx, userID.String())

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestListByUserID_DecryptError() {
	ctx := context.Background()
	userID := uuid.New()
	encKey := []byte("encrypted")
	mockModel := &llmapikey.StoredAPIKey{
		ID:           uuid.New(),
		Name:         "Test Key",
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}

	s.mockRepo.On("ListByUserID", ctx, userID.String()).Return([]*llmapikey.StoredAPIKey{
		mockModel,
	}, nil)
	s.mockEncryption.On("Decrypt", encKey).Return(nil, assert.AnError)

	result, err := s.service.ListByUserID(ctx, userID.String())

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockRepo.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestCreate_Success() {
	ctx := context.Background()
	userID := uuid.New()
	plainKey := "sk-1234567890abcdef"
	encKey := []byte("encrypted")
	createdAt := time.Now()
	params := llmapikey.CreateInput{
		Name:   "Test Key",
		Value:  plainKey,
		UserID: userID.String(),
	}
	mockModel := &llmapikey.StoredAPIKey{
		ID:           uuid.New(),
		Name:         params.Name,
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    createdAt,
	}

	s.mockEncryption.On("Encrypt", []byte(plainKey)).Return(encKey)
	s.mockRepo.On("Create", ctx, mock.MatchedBy(func(p llmapikey.CreateParams) bool {
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

func (s *ServiceTestSuite) TestCreate_RepoError() {
	ctx := context.Background()
	userID := uuid.New()
	plainKey := "sk-1234567890abcdef"
	encKey := []byte("encrypted")
	params := llmapikey.CreateInput{
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
