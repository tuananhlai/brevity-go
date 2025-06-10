package repository_test

import (
	"context"
	"testing"

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

func (s *LLMAPIKeyRepositoryTestSuite) TestCreateLLMAPIKey() {
	ctx := context.Background()
	user := s.mustCreateUser()

	err := s.llmAPIKeyRepo.Create(ctx, &model.LLMAPIKey{
		Name:         "testname",
		EncryptedKey: []byte("testencryptedkey"),
		UserID:       user.ID,
	})
	s.Require().NoError(err)
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
