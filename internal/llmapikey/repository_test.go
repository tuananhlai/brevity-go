package llmapikey_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/auth"
	"github.com/tuananhlai/brevity-go/internal/llmapikey"
	"github.com/tuananhlai/brevity-go/internal/testutil"
)

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

type RepositoryTestSuite struct {
	suite.Suite
	dbTestUtil    *testutil.DatabaseTestUtil
	authRepo      auth.Repository
	llmAPIKeyRepo llmapikey.Repository
}

func (s *RepositoryTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.authRepo = auth.NewRepository(s.dbTestUtil.DB())
	s.llmAPIKeyRepo = llmapikey.NewRepository(s.dbTestUtil.DB())
}

func (s *RepositoryTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *RepositoryTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	s.Require().NoError(err)
}

func (s *RepositoryTestSuite) TestCreateLLMAPIKey_Success() {
	ctx := context.Background()
	user := s.mustCreateUser()

	expectedName := "testname"
	expectedEncryptedKey := []byte("testencryptedkey")
	expectedUserID := user.ID

	_, err := s.llmAPIKeyRepo.Create(ctx, llmapikey.CreateParams{
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

func (s *RepositoryTestSuite) TestListByUserID_Success() {
	ctx := context.Background()
	user := s.mustCreateUser()

	expectedName := "testname"
	expectedEncryptedKey := []byte("testencryptedkey")
	expectedUserID := user.ID

	_, err := s.llmAPIKeyRepo.Create(ctx, llmapikey.CreateParams{
		Name:         expectedName,
		EncryptedKey: expectedEncryptedKey,
		UserID:       expectedUserID.String(),
	})
	s.Require().NoError(err)

	keys, err := s.llmAPIKeyRepo.ListByUserID(ctx, expectedUserID.String())
	s.Require().NoError(err)

	s.Len(keys, 1)
	s.Equal(expectedName, keys[0].Name)
	s.Equal(expectedEncryptedKey, keys[0].EncryptedKey)
	s.Equal(expectedUserID, keys[0].UserID)
}

func (s *RepositoryTestSuite) mustCreateUser() *auth.User {
	user := auth.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}

	createdUser, err := s.authRepo.CreateUser(context.Background(), user)
	s.Require().NoError(err)

	return createdUser
}
