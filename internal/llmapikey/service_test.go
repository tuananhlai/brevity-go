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
	"github.com/tuananhlai/brevity-go/internal/store"
)

func TestManager(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}

type ManagerTestSuite struct {
	suite.Suite
	manager        *llmapikey.Manager
	mockStore      *llmapikey.MockLLMApiKeyStore
	mockEncryption *llmapikey.MockCrypter
}

func (s *ManagerTestSuite) SetupTest() {
	s.mockStore = llmapikey.NewMockLLMApiKeyStore(s.T())
	s.mockEncryption = llmapikey.NewMockCrypter(s.T())
	s.manager = llmapikey.NewManager(s.mockStore, s.mockEncryption)
}

func (s *ManagerTestSuite) TestListByUserID_Success() {
	ctx := context.Background()
	userID := uuid.New()
	plainKey := "sk-1234567890abcdef"
	encKey := []byte("encrypted")
	createdAt := time.Now()
	mockModel := &store.OpenRouterAPIKey{
		ID:           uuid.New(),
		Name:         "Test Key",
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    createdAt,
	}

	s.mockStore.On("ListLLMAPIKeysByUserID", ctx, userID.String()).Return([]*store.OpenRouterAPIKey{
		mockModel,
	}, nil)
	s.mockEncryption.On("Decrypt", encKey).Return([]byte(plainKey), nil)

	result, err := s.manager.ListByUserID(ctx, userID.String())

	s.Require().NoError(err)
	s.Require().Len(result, 1)
	s.Require().Equal(mockModel.ID, result[0].ID)
	s.Require().Equal("sk-1234567", result[0].ValueFirstTen)
	s.Require().Equal("abcdef", result[0].ValueLastSix)
	s.Require().Equal(userID, result[0].UserID)
	s.Require().Equal(createdAt, result[0].CreatedAt)
	s.mockStore.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}

func (s *ManagerTestSuite) TestListByUserID_RepoError() {
	ctx := context.Background()
	userID := uuid.New()
	s.mockStore.On("ListLLMAPIKeysByUserID", ctx, userID.String()).Return(nil, assert.AnError)

	result, err := s.manager.ListByUserID(ctx, userID.String())

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ManagerTestSuite) TestListByUserID_DecryptError() {
	ctx := context.Background()
	userID := uuid.New()
	encKey := []byte("encrypted")
	mockModel := &store.OpenRouterAPIKey{
		ID:           uuid.New(),
		Name:         "Test Key",
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}

	s.mockStore.On("ListLLMAPIKeysByUserID", ctx, userID.String()).Return([]*store.OpenRouterAPIKey{
		mockModel,
	}, nil)
	s.mockEncryption.On("Decrypt", encKey).Return(nil, assert.AnError)

	result, err := s.manager.ListByUserID(ctx, userID.String())

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockStore.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}

func (s *ManagerTestSuite) TestCreate_Success() {
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
	mockModel := &store.OpenRouterAPIKey{
		ID:           uuid.New(),
		Name:         params.Name,
		EncryptedKey: encKey,
		UserID:       userID,
		CreatedAt:    createdAt,
	}

	s.mockEncryption.On("Encrypt", []byte(plainKey)).Return(encKey)
	s.mockStore.On("CreateLLMAPIKey", ctx, mock.MatchedBy(func(
		p store.CreateLLMAPIKeyParams,
	) bool {
		return p.Name == params.Name && string(p.EncryptedKey) == string(encKey) && p.UserID == userID.String()
	})).Return(mockModel, nil)

	result, err := s.manager.Create(ctx, params)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Equal(mockModel.ID, result.ID)
	s.Require().Equal(params.Name, result.Name)
	s.Require().Equal("sk-1234567", result.ValueFirstTen)
	s.Require().Equal("abcdef", result.ValueLastSix)
	s.Require().Equal(userID, result.UserID)
	s.Require().Equal(createdAt, result.CreatedAt)
	s.mockStore.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}

func (s *ManagerTestSuite) TestCreate_RepoError() {
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
	s.mockStore.On("CreateLLMAPIKey", ctx, mock.Anything).Return(nil, assert.AnError)

	result, err := s.manager.Create(ctx, params)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.mockStore.AssertExpectations(s.T())
	s.mockEncryption.AssertExpectations(s.T())
}
