package controller_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"github.com/tuananhlai/brevity-go/internal/controller"
	"github.com/tuananhlai/brevity-go/internal/store"
)

func TestArticleController(t *testing.T) {
	suite.Run(t, new(ArticleControllerTestSuite))
}

type ArticleControllerTestSuite struct {
	suite.Suite
	mockStore *controller.MockArticleStore
	router    *gin.Engine
}

func (s *ArticleControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *ArticleControllerTestSuite) BeforeTest(suiteName, testName string) {
	s.mockStore = controller.NewMockArticleStore(s.T())
	s.router = gin.Default()
	ctrl := controller.NewArticleController(s.mockStore)
	s.router.GET("/v1/article-previews", ctrl.ListPreviews)
	s.router.GET("/v1/articles/:slug", ctrl.GetBySlug)
}

func (s *ArticleControllerTestSuite) TestListPreviews_Success() {
	articleID := uuid.New()
	authorID := uuid.New()
	date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	previews := []store.ArticlePreview{
		{
			ID:                articleID,
			Slug:              "test-article",
			Title:             "Test Article",
			Description:       "Test Description",
			AuthorID:          authorID,
			AuthorUsername:    "test-author",
			AuthorDisplayName: sql.NullString{String: "Test Author", Valid: true},
			AuthorAvatarURL:   sql.NullString{String: "https://example.com/avatar.png", Valid: true},
			CreatedAt:         date,
			UpdatedAt:         date,
		},
	}
	s.mockStore.On("ListArticlesPreviews", mock.Anything).Return(previews, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/article-previews", nil)
	s.router.ServeHTTP(w, req)

	res := w.Body.String()

	s.Require().Equal(http.StatusOK, w.Code)
	s.mockStore.AssertExpectations(s.T())
	s.Require().Len(gjson.Get(res, "items").Array(), 1)
	s.Require().Equal(articleID.String(), gjson.Get(res, "items.0.id").String())
	s.Require().Equal(previews[0].Slug, gjson.Get(res, "items.0.slug").String())
	s.Require().Equal(previews[0].Title, gjson.Get(res, "items.0.title").String())
	s.Require().Equal(previews[0].Description, gjson.Get(res, "items.0.description").String())
	s.Require().Equal(authorID.String(), gjson.Get(res, "items.0.author.id").String())
	s.Require().Equal(previews[0].AuthorUsername, gjson.Get(res, "items.0.author.username").String())
	s.Require().Equal(previews[0].AuthorDisplayName.String,
		gjson.Get(res, "items.0.author.displayName").String())
	s.Require().Equal(previews[0].AuthorAvatarURL.String,
		gjson.Get(res, "items.0.author.avatarURL").String())
	s.Require().Equal(date.Format(time.RFC3339),
		gjson.Get(res, "items.0.createdAt").String())
	s.Require().Equal(date.Format(time.RFC3339),
		gjson.Get(res, "items.0.updatedAt").String())
}
