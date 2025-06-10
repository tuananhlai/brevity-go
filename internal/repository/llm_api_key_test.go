package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/testutil"
)

func TestLLMAPIKeyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LLMAPIKeyRepositoryTestSuite))
}

type LLMAPIKeyRepositoryTestSuite struct {
	suite.Suite
	dbTestUtil    *testutil.DatabaseTestUtil
	authRepo      repository.AuthRepository
	llmAPIKeyRepo repository.LLMAPIKeyRepository
}

func (s *LLMAPIKeyRepositoryTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.authRepo = repository.NewAuthRepository(s.dbTestUtil.DB())
	s.llmAPIKeyRepo = repository.NewLLMAPIKeyRepository(s.dbTestUtil.DB())
}

func (s *LLMAPIKeyRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *LLMAPIKeyRepositoryTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	s.Require().NoError(err)
}

func (s *LLMAPIKeyRepositoryTestSuite) TestCreateLLMAPIKey_Success() {
	ctx := context.Background()
	user := s.mustCreateUser()

	expectedName := "testname"
	expectedEncryptedKey := []byte("testencryptedkey")
	expectedUserID := user.ID

	err := s.llmAPIKeyRepo.Create(ctx, &model.LLMAPIKey{
		Name:         expectedName,
		EncryptedKey: expectedEncryptedKey,
		UserID:       expectedUserID,
	})
	s.Require().NoError(err)

	var actualName string
	var actualEncryptedKey []byte
	var actualUserID uuid.UUID

	row := s.dbTestUtil.DB().QueryRow("SELECT name, encrypted_key, user_id FROM llm_api_keys")
	err = row.Scan(&actualName, &actualEncryptedKey, &actualUserID)
	s.Require().NoError(err)

	s.Equal(expectedName, actualName)
	s.Equal(expectedEncryptedKey, actualEncryptedKey)
	s.Equal(expectedUserID, actualUserID)
}

func (s *LLMAPIKeyRepositoryTestSuite) TestListByUserID_Success() {
	ctx := context.Background()
	user := s.mustCreateUser()

	expectedName := "testname"
	expectedEncryptedKey := []byte("testencryptedkey")
	expectedUserID := user.ID

	err := s.llmAPIKeyRepo.Create(ctx, &model.LLMAPIKey{
		Name:         expectedName,
		EncryptedKey: expectedEncryptedKey,
		UserID:       expectedUserID,
	})
	s.Require().NoError(err)

	keys, err := s.llmAPIKeyRepo.ListByUserID(ctx, expectedUserID)
	s.Require().NoError(err)

	s.Len(keys, 1)
	s.Equal(expectedName, keys[0].Name)
	s.Equal(expectedEncryptedKey, keys[0].EncryptedKey)
	s.Equal(expectedUserID, keys[0].UserID)
}

func (s *LLMAPIKeyRepositoryTestSuite) mustCreateUser() *model.AuthUser {
	user := repository.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}

	createdUser, err := s.authRepo.CreateUser(context.Background(), user)
	s.Require().NoError(err)

	return createdUser
}
