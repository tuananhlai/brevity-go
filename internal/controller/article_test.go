package controller_test

import (
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
	"github.com/tuananhlai/brevity-go/internal/service/mocks"
)

func TestArticleController(t *testing.T) {
	suite.Run(t, new(ArticleControllerTestSuite))
}

type ArticleControllerTestSuite struct {
	suite.Suite
	mockService *mocks.ArticleService
	controller  *controller.ArticleController
	router      *gin.Engine
}

func (s *ArticleControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *ArticleControllerTestSuite) BeforeTest(suiteName, testName string) {
	s.mockService = new(mocks.ArticleService)
	s.controller = controller.NewArticleController(s.mockService)

	s.router = gin.Default()
	s.controller.RegisterRoutes(s.router)
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
			AuthorDisplayName: "Test Author",
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
	s.Require().Equal(authorID.String(), gjson.Get(res, "items.0.authorID").String())
	s.Require().Equal(previews[0].AuthorDisplayName, gjson.Get(res, "items.0.authorDisplayName").String())
	s.Require().Equal(date.Format(time.RFC3339),
		gjson.Get(res, "items.0.createdAt").String())
	s.Require().Equal(date.Format(time.RFC3339),
		gjson.Get(res, "items.0.updatedAt").String())
	s.Require().Equal(nextPageToken, gjson.Get(res, "nextPageToken").String())
}
