package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/testutil"
)

func TestDigitalAuthorRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DigitalAuthorRepositoryTestSuite))
}

type DigitalAuthorRepositoryTestSuite struct {
	suite.Suite
	dbTestUtil        *testutil.DatabaseTestUtil
	authRepo          repository.AuthRepository
	llmAPIKeyRepo     repository.LLMAPIKeyRepository
	digitalAuthorRepo repository.DigitalAuthorRepository
}

func (s *DigitalAuthorRepositoryTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.authRepo = repository.NewAuthRepository(s.dbTestUtil.DB())
	s.llmAPIKeyRepo = repository.NewLLMAPIKeyRepository(s.dbTestUtil.DB())
	s.digitalAuthorRepo = repository.NewDigitalAuthorRepository(s.dbTestUtil.DB())
}

func (s *DigitalAuthorRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *DigitalAuthorRepositoryTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	s.Require().NoError(err)
}

func (s *DigitalAuthorRepositoryTestSuite) TestCreateDigitalAuthor_Success() {
	ctx := context.Background()
	user := s.mustCreateUser()
	apiKey := s.mustCreateLLMAPIKey(user.ID.String())

	expectedDisplayName := "testdisplayname"
	expectedSystemPrompt := "testsystemprompt"
	expectedDefaultUserPrompt := "testdefaultuserprompt"
	expectedAvatarURL := "testavatarurl"

	digitalAuthor, err := s.digitalAuthorRepo.Create(ctx, repository.DigitalAuthorCreateParams{
		OwnerID:           user.ID.String(),
		DisplayName:       expectedDisplayName,
		SystemPrompt:      expectedSystemPrompt,
		DefaultUserPrompt: expectedDefaultUserPrompt,
		APIKeyID:          apiKey.ID.String(),
		AvatarURL:         expectedAvatarURL,
	})

	s.Require().NoError(err)

	err = s.dbTestUtil.DB().QueryRowContext(ctx, `
		SELECT id, owner_id, display_name, system_prompt, default_user_prompt, api_key_id, avatar_url, created_at, updated_at
		FROM digital_authors
		WHERE id = $1
	`, digitalAuthor.ID).Scan(&digitalAuthor.ID, &digitalAuthor.OwnerID,
		&digitalAuthor.DisplayName, &digitalAuthor.SystemPrompt, &digitalAuthor.DefaultUserPrompt,
		&digitalAuthor.APIKeyID, &digitalAuthor.AvatarURL, &digitalAuthor.CreatedAt, &digitalAuthor.UpdatedAt)

	s.Require().NoError(err)
	s.Equal(expectedDisplayName, digitalAuthor.DisplayName)
	s.Equal(expectedSystemPrompt, digitalAuthor.SystemPrompt)
	s.Equal(expectedDefaultUserPrompt, digitalAuthor.DefaultUserPrompt)
	s.Equal(expectedAvatarURL, digitalAuthor.AvatarURL)
	s.Equal(apiKey.ID, digitalAuthor.APIKeyID)
	s.Equal(user.ID, digitalAuthor.OwnerID)
}

func (s *DigitalAuthorRepositoryTestSuite) mustCreateUser() *model.AuthUser {
	user := repository.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}

	createdUser, err := s.authRepo.CreateUser(context.Background(), user)
	s.Require().NoError(err)

	return createdUser
}

func (s *DigitalAuthorRepositoryTestSuite) mustCreateLLMAPIKey(userID string) *model.LLMAPIKey {
	apiKey := repository.LLMAPIKeyCreateParams{
		Name:         "testapikey",
		EncryptedKey: []byte("encryptedKey"),
		UserID:       userID,
	}

	createdAPIKey, err := s.llmAPIKeyRepo.Create(context.Background(), apiKey)
	s.Require().NoError(err)

	return createdAPIKey
}
