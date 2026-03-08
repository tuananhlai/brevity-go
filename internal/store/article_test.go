package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/tuananhlai/brevity-go/internal/store"
	"github.com/tuananhlai/brevity-go/internal/testutil"
)

func TestArticleStore(t *testing.T) {
	suite.Run(t, new(ArticleStoreTestSuite))
}

type ArticleStoreTestSuite struct {
	suite.Suite
	dbTestUtil *testutil.DatabaseTestUtil
	store      *store.Store
}

func (s *ArticleStoreTestSuite) SetupSuite() {
	var err error
	s.dbTestUtil, err = testutil.NewDatabaseTestUtil()
	s.Require().NoError(err)

	s.store = store.New(s.dbTestUtil.DB())
}

func (s *ArticleStoreTestSuite) BeforeTest(suiteName, testName string) {
	err := s.dbTestUtil.Reset()
	s.Require().NoError(err)
}

func (s *ArticleStoreTestSuite) TearDownSuite() {
	err := s.dbTestUtil.Teardown()
	s.Require().NoError(err)
}

func (s *ArticleStoreTestSuite) TestCreateArticle_Success() {
	ctx := context.Background()
	author := s.mustCreateUser()

	article := &store.Article{
		Title:    "Test Article",
		Content:  "This is a test article",
		AuthorID: author.ID,
	}
	err := s.store.CreateArticle(ctx, article)
	s.Require().NoError(err)
}

func (s *ArticleStoreTestSuite) TestListArticlesPreviews_Success() {
	ctx := context.Background()
	author := s.mustCreateUser()
	newArticle := s.mustCreateArticle(author.ID)

	previews, err := s.store.ListArticlesPreviews(ctx)

	s.Require().NoError(err)
	s.Require().Len(previews, 1)
	s.Require().Equal(newArticle.Title, previews[0].Title)
	s.Require().Equal(newArticle.Description, previews[0].Description)
	s.Require().Equal(newArticle.AuthorID, previews[0].AuthorID)
}

func (s *ArticleStoreTestSuite) TestGetArticleBySlug_Success() {
	ctx := context.Background()
	author := s.mustCreateUser()
	newArticle := s.mustCreateArticle(author.ID)

	article, err := s.store.GetArticleBySlug(ctx, newArticle.Slug)

	s.Require().NoError(err)
	s.Require().Equal(newArticle.Slug, article.Slug)
	s.Require().Equal(newArticle.Title, article.Title)
	s.Require().Equal(newArticle.AuthorID, article.AuthorID)
	s.Require().WithinDuration(newArticle.CreatedAt, article.CreatedAt, time.Microsecond)
	s.Require().WithinDuration(newArticle.UpdatedAt, article.UpdatedAt, time.Microsecond)
	s.Require().Equal(newArticle.Content, article.Content)
	s.Require().Equal(author.Username, article.AuthorUsername)
}

func (s *ArticleStoreTestSuite) mustCreateUser() *store.User {
	user := store.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}

	createdUser, err := s.store.CreateUser(context.Background(), user)
	s.Require().NoError(err)
	s.mustCreateDigitalAuthor(createdUser.ID)

	return createdUser
}

func (s *ArticleStoreTestSuite) mustCreateDigitalAuthor(id uuid.UUID) {
	_, err := s.dbTestUtil.DB().ExecContext(context.Background(), `
		INSERT INTO digital_authors (id, display_name, system_prompt)
		VALUES ($1, $2, $3)
	`, id, "Test Author Bot", "Write helpful articles")
	s.Require().NoError(err)
}

func (s *ArticleStoreTestSuite) mustCreateArticle(authorID uuid.UUID) *store.Article {
	article := &store.Article{
		Title:    "Test Article",
		Content:  "This is a test article",
		AuthorID: authorID,
	}

	err := s.store.CreateArticle(context.Background(), article)
	s.Require().NoError(err)

	err = s.dbTestUtil.DB().GetContext(context.Background(), article, `
		SELECT * FROM articles LIMIT 1`,
	)
	s.Require().NoError(err)

	return article
}
