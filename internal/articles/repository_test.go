package articles_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/articles"
	"github.com/tuananhlai/brevity-go/internal/auth"
	"github.com/tuananhlai/brevity-go/internal/testutil"
)

func TestArticleRepository(t *testing.T) {
	suite.Run(t, new(ArticleRepositoryTestSuite))
}

type ArticleRepositoryTestSuite struct {
	suite.Suite
	dbTestUtil  *testutil.DatabaseTestUtil
	authRepo    auth.Repository
	articleRepo articles.Repository
}

func (s *ArticleRepositoryTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.authRepo = auth.NewRepository(s.dbTestUtil.DB())
	s.articleRepo = articles.NewRepository(s.dbTestUtil.DB())
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
	author := s.mustCreateUser()

	article := &articles.Article{
		Title:    "Test Article",
		Content:  "This is a test article",
		AuthorID: author.ID,
	}
	err := s.articleRepo.Create(ctx, article)
	s.Require().NoError(err)
}

func (s *ArticleRepositoryTestSuite) TestListPreviews_Success() {
	ctx := context.Background()
	author := s.mustCreateUser()
	newArticle := s.mustCreateArticle(author.ID)

	previews, _, err := s.articleRepo.ListPreviews(ctx, 100)

	s.Require().NoError(err)
	s.Require().Len(previews, 1)
	s.Require().Equal(newArticle.Title, previews[0].Title)
	s.Require().Equal(newArticle.Description, previews[0].Description)
	s.Require().Equal(newArticle.AuthorID, previews[0].AuthorID)
}

func (s *ArticleRepositoryTestSuite) TestGetBySlug_Success() {
	ctx := context.Background()
	author := s.mustCreateUser()
	newArticle := s.mustCreateArticle(author.ID)

	article, err := s.articleRepo.GetBySlug(ctx, newArticle.Slug)

	s.Require().NoError(err)
	s.Require().Equal(newArticle.Slug, article.Slug)
	s.Require().Equal(newArticle.Title, article.Title)
	s.Require().Equal(newArticle.AuthorID, article.AuthorID)
	s.Require().WithinDuration(newArticle.CreatedAt, article.CreatedAt, time.Microsecond)
	s.Require().WithinDuration(newArticle.UpdatedAt, article.UpdatedAt, time.Microsecond)
	s.Require().Equal(newArticle.Content, article.Content)
	s.Require().Equal(author.Username, article.AuthorUsername)
}

func (s *ArticleRepositoryTestSuite) mustCreateUser() *auth.User {
	user := auth.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}

	createdUser, err := s.authRepo.CreateUser(context.Background(), user)
	s.Require().NoError(err)

	return createdUser
}

func (s *ArticleRepositoryTestSuite) mustCreateArticle(authorID uuid.UUID) *articles.Article {
	article := &articles.Article{
		Title:    "Test Article",
		Content:  "This is a test article",
		AuthorID: authorID,
	}

	err := s.articleRepo.Create(context.Background(), article)
	s.Require().NoError(err)

	err = s.dbTestUtil.DB().GetContext(context.Background(), article, `
		SELECT * FROM articles LIMIT 1`,
	)
	s.Require().NoError(err)

	return article
}
