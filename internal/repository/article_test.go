package repository_test

import (
	"context"
	"testing"

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
	user := repository.CreateUserParams{
		Username:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: []byte("passwordHash"),
	}
	createdUser, err := s.authRepo.CreateUser(ctx, user)
	s.Require().NoError(err)

	article := &model.Article{
		Title:    "Test Article",
		Content:  "This is a test article",
		AuthorID: createdUser.ID,
	}
	err = s.articleRepo.Create(context.Background(), article)
	s.Require().NoError(err)
}
