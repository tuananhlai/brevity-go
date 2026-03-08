package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/store"
	"github.com/tuananhlai/brevity-go/internal/testutil"
)

func TestLLMAPIKeyStore(t *testing.T) {
	suite.Run(t, new(LLMAPIKeyStoreTestSuite))
}

type LLMAPIKeyStoreTestSuite struct {
	suite.Suite
	dbTestUtil *testutil.DatabaseTestUtil
	store      *store.Store
}

func (s *LLMAPIKeyStoreTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.store = store.New(s.dbTestUtil.DB())
}

func (s *LLMAPIKeyStoreTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *LLMAPIKeyStoreTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	s.Require().NoError(err)
}

func (s *LLMAPIKeyStoreTestSuite) TestCreateLLMAPIKey_Success() {
	ctx := context.Background()
	user := s.mustCreateUser()

	expectedName := "testname"
	expectedEncryptedKey := []byte("testencryptedkey")
	expectedUserID := user.ID

	_, err := s.store.CreateLLMAPIKey(ctx, store.CreateLLMAPIKeyParams{
		Name:         expectedName,
		EncryptedKey: expectedEncryptedKey,
		UserID:       expectedUserID.String(),
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

func (s *LLMAPIKeyStoreTestSuite) TestListLLMAPIKeysByUserID_Success() {
	ctx := context.Background()
	user := s.mustCreateUser()

	expectedName := "testname"
	expectedEncryptedKey := []byte("testencryptedkey")
	expectedUserID := user.ID

	_, err := s.store.CreateLLMAPIKey(ctx, store.CreateLLMAPIKeyParams{
		Name:         expectedName,
		EncryptedKey: expectedEncryptedKey,
		UserID:       expectedUserID.String(),
	})
	s.Require().NoError(err)

	keys, err := s.store.ListLLMAPIKeysByUserID(ctx, expectedUserID.String())
	s.Require().NoError(err)

	s.Len(keys, 1)
	s.Equal(expectedName, keys[0].Name)
	s.Equal(expectedEncryptedKey, keys[0].EncryptedKey)
	s.Equal(expectedUserID, keys[0].UserID)
}

func (s *LLMAPIKeyStoreTestSuite) mustCreateUser() *store.User {
	user := store.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}

	createdUser, err := s.store.CreateUser(context.Background(), user)
	s.Require().NoError(err)

	return createdUser
}
