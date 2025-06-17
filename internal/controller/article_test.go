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
	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/service"
)

func TestArticleController(t *testing.T) {
	suite.Run(t, new(ArticleControllerTestSuite))
}

type ArticleControllerTestSuite struct {
	suite.Suite
	mockService *service.MockArticleService
	router      *gin.Engine
}

func (s *ArticleControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *ArticleControllerTestSuite) BeforeTest(suiteName, testName string) {
	s.mockService = service.NewMockArticleService(s.T())
	s.router = gin.Default()
	controller := controller.NewArticleController(s.mockService)
	s.router.GET("/v1/article-previews", controller.ListPreviews)
	s.router.GET("/v1/articles/:slug", controller.GetBySlug)
}

func (s *ArticleControllerTestSuite) TestListPreviews_Success() {
	articleID := uuid.New()
	authorID := uuid.New()
	date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	previews := []model.ArticlePreview{
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
	nextPageToken := "next-token"
	s.mockService.On("ListPreviews", mock.Anything, 50, mock.Anything).Return(previews, nextPageToken, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/article-previews?pageSize=50", nil)
	s.router.ServeHTTP(w, req)

	res := w.Body.String()

	s.Require().Equal(http.StatusOK, w.Code)
	s.mockService.AssertExpectations(s.T())
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
	s.Require().Equal(nextPageToken, gjson.Get(res, "nextPageToken").String())
}
