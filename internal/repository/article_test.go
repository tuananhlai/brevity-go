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

func TestArticleRepository(t *testing.T) {
	suite.Run(t, new(ArticleRepositoryTestSuite))
}

type ArticleRepositoryTestSuite struct {
	suite.Suite
	dbTestUtil  *testutil.DatabaseTestUtil
	authRepo    repository.AuthRepository
	articleRepo repository.ArticleRepository
}

func (s *ArticleRepositoryTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.authRepo = repository.NewAuthRepository(s.dbTestUtil.DB())
	s.articleRepo = repository.NewArticleRepository(s.dbTestUtil.DB())
}

func (s *ArticleRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *ArticleRepositoryTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	s.Require().NoError(err)
}

func (s *ArticleRepositoryTestSuite) TestCreateArticle_Success() {
	ctx := context.Background()
	authorID := mustCreateUser(s)

	article := &model.Article{
		Title:    "Test Article",
		Content:  "This is a test article",
		AuthorID: authorID,
	}
	err := s.articleRepo.Create(ctx, article)
	s.Require().NoError(err)
}

func (s *ArticleRepositoryTestSuite) TestListPreviews_Success() {
	ctx := context.Background()
	authorID := mustCreateUser(s)
	newArticle := mustCreateArticle(s, authorID)

	previews, err := s.articleRepo.ListPreviews(ctx)

	s.Require().NoError(err)
	s.Require().Len(previews, 1)
	s.Require().Equal(newArticle.Title, previews[0].Title)
	s.Require().Equal(newArticle.Description, previews[0].Description)
	s.Require().Equal(newArticle.AuthorID, previews[0].AuthorID)
}

func mustCreateUser(s *ArticleRepositoryTestSuite) uuid.UUID {
	user := repository.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}

	createdUser, err := s.authRepo.CreateUser(context.Background(), user)
	s.Require().NoError(err)

	return createdUser.ID
}

func mustCreateArticle(s *ArticleRepositoryTestSuite, authorID uuid.UUID) *model.Article {
	article := &model.Article{
		Title:    "Test Article",
		Content:  "This is a test article",
		AuthorID: authorID,
	}

	err := s.articleRepo.Create(context.Background(), article)
	s.Require().NoError(err)

	return article
}
